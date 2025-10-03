package ocr

import (
	"context"
)

func (r Repository) ScanBytes(ctx context.Context, file []byte) (string, error) {
	return "", nil
	//client := gosseract.NewClient()
	//defer client.Close()
	//
	//err := client.SetLanguage("eng+rus")
	//if err != nil {
	//	return "", err
	//}
	//
	//err = client.SetImageFromBytes(file)
	//if err != nil {
	//	return "", err
	//}
	//
	//err = client.SetPageSegMode(gosseract.PSM_AUTO) // Автоопределение блоков текста
	//if err != nil {
	//	return "", err
	//}
	//
	//err = client.SetVariable("user_defined_dpi", "300") // Увеличение точности
	//if err != nil {
	//	return "", err
	//}

	//client.SetVariable("tessedit_char_whitelist", "0123456789.,mg/dL%") // Ограничение символов
	//client.SetVariable("preserve_interword_spaces", "1") // Сохранение пробелов между словами
	//client.SetVariable("load_system_dawg", "F") // Отключение системных словарей для уменьшения ошибок
	//client.SetVariable("load_freq_dawg", "F")   // Отключение частотных словарей

	//client.SetVariable("ocr_engine_mode", "1")          // Используем только LSTM (без старого OCR)
	//client.SetVariable("tessedit_ocr_engine_mode", "1") // Отключаем гибридный режим (улучшает точность)

	//client.SetVariable("tessedit_char_whitelist", "0123456789.,mg/dL%АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯabcdefghijklmnopqrstuvwxyz")
	//client.SetVariable("load_system_dawg", "F")          // -c load_system_dawg=false
	//client.SetVariable("load_freq_dawg", "F")            // -c load_freq_dawg=false
	//client.SetVariable("preserve_interword_spaces", "1") // -c preserve_interword_spaces=1

	//return client.Text()
}

// preprocessImage принимает изображение в байтах, обрабатывает его и возвращает обработанное изображение в байтах
//func preprocessImage(inputBytes []byte) ([]byte, error) {
//	// Декодируем изображение
//	img, err := gocv.IMDecode(inputBytes, gocv.IMReadGrayScale)
//	if err != nil || img.Empty() {
//		return nil, fmt.Errorf("ошибка декодирования изображения: %v", err)
//	}
//	defer img.Close()
//
//	// Используем CLAHE (улучшение контраста без потери деталей)
//	clahe := gocv.NewCLAHE()
//	defer clahe.Close()
//	clahe.Apply(img, &img)
//
//	// Адаптивная бинаризация (избегает жесткой обрезки деталей)
//	gocv.AdaptiveThreshold(img, &img, 255, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, 11, 2)
//
//	// Кодируем обратно в байты
//	buf := new(bytes.Buffer)
//	mat, err := gocv.IMEncode(gocv.PNGFileExt, img) // Используем PNG для сохранения качества
//	if err != nil {
//		return nil, fmt.Errorf("ошибка кодирования изображения: %v", err)
//	}
//	defer mat.Close()
//
//	_, err = buf.Write(mat.GetBytes())
//	if err != nil {
//		return nil, fmt.Errorf("ошибка записи обработанного изображения в буфер: %v", err)
//	}
//
//	return buf.Bytes(), nil
//}
