import json

# 📌 Исходный текст (OCR-распознанный)
text = """
Пациент: Иванов Иван Иванович
Дата анализа: 23.02.2025
Лаборатория: Invitro

Гемоглобин 140 г/л 120-158
Лейкоциты 6.2 10⁹/л 4.00-10.50
Глюкоза 5.6 ммоль/л 4.1-5.9
"""

# 📌 Разметка для обучения BERT (IOB-формат)
annotations = {
    "Иванов Иван Иванович": "LABEL_PATIENT",
    "23.02.2025": "LABEL_DATE",
    "Invitro": "LABEL_LAB",
    "Гемоглобин": "LABEL_ANALYSIS",
    "140": "LABEL_VALUE",
    "г/л": "LABEL_UNIT",
    "120-158": "LABEL_REF",
    "Лейкоциты": "LABEL_ANALYSIS",
    "6.2": "LABEL_VALUE",
    "10⁹/л": "LABEL_UNIT",
    "4.00-10.50": "LABEL_REF",
    "Глюкоза": "LABEL_ANALYSIS",
    "5.6": "LABEL_VALUE",
    "ммоль/л": "LABEL_UNIT",
    "4.1-5.9": "LABEL_REF",
}

# 📌 Генерируем разметку в формате IOB
tokenized_data = []
for line in text.split("\n"):
    for word in line.split():
        label = annotations.get(word, "O")
        prefix = "B-" if label != "O" else "O"
        tokenized_data.append(f"{word} {prefix}{label}")

# 📌 Записываем в файл
with open("train_data.txt", "w", encoding="utf-8") as f:
    f.write("\n".join(tokenized_data))

print("✅ Размеченный файл train_data.txt создан!")
