import os
import subprocess

# 📌 Папка с текстовыми файлами
TEXT_FILES_DIR = "./texts"  # Укажите путь к папке с текстами
OUTPUT_DIR = "./images"     # Папка для сохранения изображений
FONT_NAME = "Arial"         # Используемый шрифт
FONTS_DIR = "/usr/share/fonts"  # Каталог шрифтов

# 📌 Убедимся, что папка для изображений существует
os.makedirs(OUTPUT_DIR, exist_ok=True)

# 📌 Проходим по всем текстовым файлам в папке
for filename in os.listdir(TEXT_FILES_DIR):
    if filename.endswith(".txt"):
        text_file_path = os.path.join(TEXT_FILES_DIR, filename)
        output_base = os.path.join(OUTPUT_DIR, os.path.splitext(filename)[0])

        # 📌 Запускаем text2image
        command = [
            "text2image",
            "--text", text_file_path,
            "--outputbase", output_base,
            "--font", FONT_NAME,
            "--fonts_dir", FONTS_DIR
        ]

        print(f"🖼 Генерируем изображение для: {filename}")
        subprocess.run(command, check=True)

print("✅ Генерация изображений завершена!")
