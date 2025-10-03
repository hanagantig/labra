# Create a robust scraper the user can run locally to parse all biomarkers from the given Lab4U page
# It handles:
# - Static HTML product cards
# - Embedded JSON (window.__INITIAL_STATE__ / application/ld+json)
# - Pagination within the "frequently_search" section
# Outputs a CSV with columns: name, url, section, city, price_current, price_old, unit, biomaterial (best-effort)
#
# Usage:
#   python lab4u_frequently_search_scraper.py --city moscow --out biomarkers_lab4u.csv
#
# Requirements:
#   pip install requests beautifulsoup4 lxml

script = r'''#!/usr/bin/env python3
import re
import csv
import json
import time
import argparse
from urllib.parse import urljoin, urlparse
import requests
from bs4 import BeautifulSoup

BASE = "https://lab4u.ru"

def fetch(url, session, retries=3):
    for i in range(retries):
        try:
            resp = session.get(url, timeout=20, headers={
                "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0 Safari/537.36"
            })
            if resp.status_code == 200:
                return resp.text
        except requests.RequestException:
            pass
        time.sleep(1 + i)
    raise RuntimeError(f"Failed to fetch {url}")

def parse_cards(html, base=BASE):
    soup = BeautifulSoup(html, "lxml")
    cards = []

    # 1) Try to parse product tiles
    for card in soup.select("[data-qa='catalog-item'], .catalog__item, .product-card, .CatalogCard_root__*, .products__item"):
        name_el = card.select_one("a, [data-qa='catalog-item-title'], .product-card__title, .card__title")
        name = name_el.get_text(strip=True) if name_el else None
        href = name_el.get("href") if name_el and name_el.has_attr("href") else None
        url = urljoin(base, href) if href else None

        price_now = None
        price_old = None
        for sel in [".price__current", ".price__new", "[data-qa='price-current']"]:
            p = card.select_one(sel)
            if p:
                price_now = re.sub(r"[^\d,\.]", "", p.get_text(" ", strip=True))
                break
        for sel in [".price__old", "[data-qa='price-old']"]:
            p = card.select_one(sel)
            if p:
                price_old = re.sub(r"[^\d,\.]", "", p.get_text(" ", strip=True))
                break

        unit = None
        biomaterial = None
        # heuristic text scraping
        info_text = card.get_text(" ", strip=True).lower()
        if "сыворотк" in info_text:
            biomaterial = "Serum"
        elif "плазм" in info_text:
            biomaterial = "Plasma"
        elif "кров" in info_text:
            biomaterial = "Blood"
        elif "моч" in info_text:
            biomaterial = "Urine"
        elif "кал" in info_text:
            biomaterial = "Feces"
        elif "слюн" in info_text:
            biomaterial = "Saliva"

        if name:
            cards.append({
                "name": name,
                "url": url,
                "price_current": price_now,
                "price_old": price_old,
                "unit": unit,
                "biomaterial": biomaterial
            })

    # 2) Try embedded JSON in scripts
    text = soup.get_text("\n", strip=False)
    m = re.search(r"window\.__INITIAL_STATE__\s*=\s*(\{.*?\});", html, flags=re.S)
    if m:
        try:
            state = json.loads(m.group(1))
            # heuristic path to products
            for path in [
                ["catalog", "items"],
                ["page", "items"],
                ["items"],
            ]:
                d = state
                ok = True
                for p in path:
                    if isinstance(d, dict) and p in d:
                        d = d[p]
                    else:
                        ok = False
                        break
                if ok and isinstance(d, list):
                    for it in d:
                        title = it.get("title") or it.get("name")
                        href = it.get("url") or it.get("link")
                        price = it.get("price") or it.get("cost") or it.get("currentPrice")
                        if title:
                            cards.append({
                                "name": title,
                                "url": urljoin(base, href) if href else None,
                                "price_current": price,
                                "price_old": None,
                                "unit": it.get("unit"),
                                "biomaterial": it.get("biomaterial") or it.get("biomaterialName")
                            })
        except Exception:
            pass

    # 3) JSON-LD
    for tag in soup.select("script[type='application/ld+json']"):
        try:
            data = json.loads(tag.string)
        except Exception:
            continue
        if isinstance(data, dict):
            data = [data]
        for obj in data:
            if isinstance(obj, dict) and obj.get("@type") in ("Product", "Offer"):
                name = obj.get("name")
                url = obj.get("url")
                offer = obj.get("offers") or {}
                price_now = None
                if isinstance(offer, dict):
                    price_now = offer.get("price")
                if name:
                    cards.append({
                        "name": name,
                        "url": url,
                        "price_current": price_now,
                        "price_old": None,
                        "unit": None,
                        "biomaterial": None
                    })

    return cards

def next_page_url(html, current_url):
    soup = BeautifulSoup(html, "lxml")
    for a in soup.select("a[rel='next'], .pagination__next a, a:contains('Следующая')"):
        href = a.get("href")
        if href:
            return urljoin(current_url, href)
    # Sometimes pagination via numbered links
    for a in soup.select(".pagination a"):
        if ">" in a.get_text(strip=True) or "След" in a.get_text(strip=True):
            href = a.get("href")
            if href:
                return urljoin(current_url, href)
    return None

def scrape_section(city, out_csv):
    session = requests.Session()
    start_url = f"{BASE}/{city}/store/section/frequently_search/"
    url = start_url
    seen = set()
    all_rows = []
    while url and url not in seen:
        seen.add(url)
        html = fetch(url, session)
        rows = parse_cards(html, base=BASE)
        all_rows.extend(rows)
        url = next_page_url(html, url)

    # Deduplicate by name+url
    uniq = {}
    for r in all_rows:
        key = (r.get("name"), r.get("url"))
        if key not in uniq:
            uniq[key] = r
    rows = list(uniq.values())

    with open(out_csv, "w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=["name","url","price_current","price_old","unit","biomaterial","section","city"])
        w.writeheader()
        for r in rows:
            r["section"] = "frequently_search"
            r["city"] = city
            w.writerow(r)

    print(f"Wrote {len(rows)} rows to {out_csv}")

if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("--city", default="moscow", help="City segment in URL (e.g., moscow)")
    ap.add_argument("--out", default="biomarkers_lab4u.csv", help="Output CSV path")
    args = ap.parse_args()
    scrape_section(args.city, args.out)
'''
path = "./lab4u_frequently_search_scraper.py"
with open(path, "w", encoding="utf-8") as f:
    f.write(script)

path
