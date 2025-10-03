#!/bin/bash

# 📌 Папка с текстовыми файлами
TEXT_FILES_DIR="./texts"  # Укажите путь к папке с файлами
OUTPUT_DIR="./images"     # Папка для сохранения изображений
FONT_NAME="Arial"         # Название шрифта
FONTS_DIR="/usr/share/fonts"  # Каталог шрифтов

# 📌 Убедимся, что папка для изображений существует
mkdir -p "$OUTPUT_DIR"

# 📌 Проходим по всем текстовым файлам в папке
for file in "$TEXT_FILES_DIR"/*.txt; do
    if [[ -f "$file" ]]; then
        filename=$(basename -- "$file")
        output_base="$OUTPUT_DIR/${filename%.gt.txt}"

        echo "🖼 Генерируем изображение для: $filename"

        # 📌 Запуск text2image
        text2image --text "$file" \
                   --outputbase "$output_base" \
                   --font "$FONT_NAME" \
                   --fonts_dir "$FONTS_DIR"

        tesseract "${output_base}.tif" $output_base --psm 6 lstm.train


        echo "✅ Изображение сохранено в $OUTPUT_DIR"
    fi
done

echo "🎉 Генерация завершена!"
