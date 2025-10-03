from transformers import pipeline
import json

# 📌 Загружаем предобученную модель BERT для медицинского NER
ner_pipeline = pipeline("ner", model="d4data/biomedical-ner-all", aggregation_strategy="simple")

VALID_LABELS = {"B-Test", "I-Test", "B-Disease", "I-Disease", "B-Measurement", "I-Measurement"}

def extract_medical_data(text):
    """Функция для извлечения данных из текста анализов"""
    entities = ner_pipeline(text)

    extracted_data = {"tests": [], "diseases": [], "measurements": []}

    for entity in entities:
        word, label = entity["word"], entity["entity_group"]

        print(label, word)
        if label in {"B-Test", "I-Test"}:
            extracted_data["tests"].append(word)
        elif label in {"B-Disease", "I-Disease"}:
            extracted_data["diseases"].append(word)
        elif label in {"B-Measurement", "I-Measurement"}:
            extracted_data["measurements"].append(word)

    return extracted_data

if __name__ == "__main__":
    # 📌 Ввод текста (можно заменить на OCR-распознанный текст)
    text = """
    Пациент: Иванов Иван Иванович
    Дата анализа: 23.02.2025
    Лаборатория: Invitro

    Гемоглобин 140 г/л 120-158
    Лейкоциты 6.2 10⁹/л 4.00-10.50
    Глюкоза 5.6 ммоль/л 4.1-5.9
    """

    # 📌 Обрабатываем текст через BERT
    structured_data = extract_medical_data(text)

    # 📌 Выводим JSON с результатами
    print(json.dumps(structured_data, indent=2, ensure_ascii=False))
