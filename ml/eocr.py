import easyocr

# Создаём объект распознавателя
reader = easyocr.Reader(['ru'])  # Поддержка русского и английского языка

# Задаём путь к изображению
image_path = "./001.jpg"

# Распознаём текст
results = reader.readtext(image_path)

# Выводим результаты
for (bbox, text, prob) in results:
    print(text)
#     print(f"Текст: {text}, Достоверность: {prob:.4f}")
