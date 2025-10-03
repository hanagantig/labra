package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"labra/internal/entity"
	"net/http"
	"net/url"
)

func (r Repository) Completions(ctx context.Context, message string) (entity.CheckupResults, error) {
	apiURL, err := url.JoinPath(baseURL, "chat/completions")
	if err != nil {
		return entity.CheckupResults{}, err
	}

	reqSystemMsg := Message{
		Role:    "system",
		Content: "The user will provide some medical test report. Please parse the report, find 'patient name', 'patient birthdate', 'markers', 'marker unit' and 'marker value'. If you find mistakes - correct them and output them in JSON format. EXAMPLE INPUT: Ист. бол. /Амб. карта: Ф.И.О.: Иванов Иван Иванович Дата рождения: 24/11/2003 Пол: Женский Дата/время взятия материала: 31/01/2024 09:30 Дата доставки материала: 01/02/2024 Номер заказа: 977939198001 Номер образца: 977939198001 Страховая компания: Стр. полис: Серия Номер Гематология Наименование теста          Результат               Единицы измерения Общий анализ крови (СВС/Diff) с лейкоцитарной формулой Гемоглобин  118 г/л 120-158 Эритроциты 4.09 10^12/л 3.90-5.20 Средний объём эритроцитов 87.3 фл 81.0-100.0 Среднее содержание Hb в эритроците 28,9 пг 26.0-34.0 Средняя концентрация НЬ в эритроците 331 г/л 310-370 EXAMPLE JSON OUTPUT: { 'patient_name': 'Иванов Иван Иванович', 'patient_bdate':'24/11/2003', 'markers': [{'name':'Гемоглобин', 'unit':'г/л', 'value':118}, {'name':'Эритроциты', 'unit':'10^12/л', 'value': 4.09}, {'name':'Средний объём эритроцитов', 'unit':'фл', 'value': 87.3}]",
	}

	requestBody := Request{
		Model:    "deepseek-chat",
		Messages: []Message{reqSystemMsg, {Role: "user", Content: message}},
		Stream:   false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return entity.CheckupResults{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.apiKey))

	//resp, err := r.client.Do(req)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}

	//body, err := io.ReadAll(resp.Body)

	//cnt := `json\n{\n  \"patient_name\": \"Жанатян Мария Мартиновна\",\n  \"patient_bdate\": \"24/11/2003\",\n  \"markers\": [\n    {\"name\": \"Гемоглобин\", \"unit\": \"г/л\", \"value\": 118},\n    {\"name\": \"Эритроциты\", \"unit\": \"10^12/л\", \"value\": 4.09},\n    {\"name\": \"Средний объём эритроцитов\", \"unit\": \"фл\", \"value\": 81.3},\n    {\"name\": \"Среднее содержание Hb в эритроците\", \"unit\": \"пг\", \"value\": 28.79},\n    {\"name\": \"Средняя концентрация Hb в эритроците\", \"unit\": \"г/л\", \"value\": 331},\n    {\"name\": \"Гетерогенность эритроцитов по объёму\", \"unit\": \"%\", \"value\": 12.7},\n    {\"name\": \"Гематокрит\", \"unit\": \"%\", \"value\": 38.7},\n    {\"name\": \"Тромбоциты\", \"unit\": \"10^9/л\", \"value\": 189},\n    {\"name\": \"Средний объём тромбоцитов\", \"unit\": \"фл\", \"value\": 10.9},\n    {\"name\": \"Гетерогенность тромбоцитов по объёму\", \"unit\": \"\", \"value\": 16.56},\n    {\"name\": \"Тромбокрит\", \"unit\": \"\", \"value\": 0.11270736},\n    {\"name\": \"Лейкоциты\", \"unit\": \"10^9/л\", \"value\": 4.95},\n    {\"name\": \"Нейтрофилы\", \"unit\": \"*\", \"value\": 4430},\n    {\"name\": \"Эозинофилы\", \"unit\": \"%\", \"value\": 1.60},\n    {\"name\": \"Базофилы\", \"unit\": \"*\", \"value\": 0.80},\n    {\"name\": \"Моноциты\", \"unit\": \"%\", \"value\": 5.00},\n    {\"name\": \"Лимфоциты\", \"unit\": \"*\", \"value\": 48.30},\n    {\"name\": \"Нейтрофилы\", \"unit\": \"10^9/л\", \"value\": 1.75},\n    {\"name\": \"Эозинофилы\", \"unit\": \"10^9/л\", \"value\": 0.06},\n    {\"name\": \"Базофилы\", \"unit\": \"10^9/л\", \"value\": 0.03},\n    {\"name\": \"Моноциты\", \"unit\": \"10^9/л\", \"value\": 0.20},\n    {\"name\": \"Лимфоциты\", \"unit\": \"10^9/л\", \"value\": 1.91},\n    {\"name\": \"СОЭ (Вестергрен)\", \"unit\": \"мм/час\", \"value\": 2}\n  ]\n}\n`
	//che1 := `"id":"bf513443-2ee2-43e3-9da5-2632608dbea1","object":"chat.completion","created":1741117045,"model":"deepseek-chat","choices":[{"index":0,"message":{"role":"assistant","content":"`
	//che2 := `"},"logprobs":null,"finish_reason":"stop"}],"usage":{"prompt_tokens":1199,"completion_tokens":685,"total_tokens":1884,"prompt_tokens_details":{"cached_tokens":1152},"prompt_cache_hit_tokens":1152,"prompt_cache_miss_tokens":47},"system_fingerprint":"fp_3a5770e1b4_prod0225"`
	//body := []byte(che1 + cnt + che2)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}
	//defer resp.Body.Close()
	//
	//var response DeepSeekResponse
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	return entity.CheckupResults{}, err
	//}
	//
	var parsedContent ParsedContent
	//content := response.Choices[0].Message.Content
	//content = strings.TrimPrefix(content, "```json")
	//content = strings.TrimSuffix(content, "```")
	//content = strings.TrimSpace(content)
	content := "{\n  \"patient_name\": \"Жанатян Мария Мартиновна\",\n  \"patient_bdate\": \"24/11/2003\",\n  \"markers\": [\n    {\"name\": \"Гемоглобин\", \"unit\": \"г/л\", \"value\": 118},\n    {\"name\": \"Эритроциты\", \"unit\": \"10^12/л\", \"value\": 4.09},\n    {\"name\": \"Средний объём эритроцитов\", \"unit\": \"фл\", \"value\": 81.3},\n    {\"name\": \"Среднее содержание Hb в эритроците\", \"unit\": \"пг\", \"value\": 28.79},\n    {\"name\": \"Средняя концентрация Hb в эритроците\", \"unit\": \"г/л\", \"value\": 331},\n    {\"name\": \"Гетерогенность эритроцитов по объёму\", \"unit\": \"%\", \"value\": 12.7},\n    {\"name\": \"Гематокрит\", \"unit\": \"%\", \"value\": 38.7},\n    {\"name\": \"Тромбоциты\", \"unit\": \"10^9/л\", \"value\": 189},\n    {\"name\": \"Средний объём тромбоцитов\", \"unit\": \"фл\", \"value\": 10.9},\n    {\"name\": \"Гетерогенность тромбоцитов по объёму\", \"unit\": \"\", \"value\": 16.56},\n    {\"name\": \"Тромбокрит\", \"unit\": \"\", \"value\": 0.11270736},\n    {\"name\": \"Лейкоциты\", \"unit\": \"10^9/л\", \"value\": 4.95},\n    {\"name\": \"Нейтрофилы\", \"unit\": \"*\", \"value\": 4430},\n    {\"name\": \"Эозинофилы\", \"unit\": \"%\", \"value\": 1.60},\n    {\"name\": \"Базофилы\", \"unit\": \"*\", \"value\": 0.80},\n    {\"name\": \"Моноциты\", \"unit\": \"%\", \"value\": 5.00},\n    {\"name\": \"Лимфоциты\", \"unit\": \"*\", \"value\": 48.30},\n    {\"name\": \"Нейтрофилы\", \"unit\": \"10^9/л\", \"value\": 1.75},\n    {\"name\": \"Эозинофилы\", \"unit\": \"10^9/л\", \"value\": 0.06},\n    {\"name\": \"Базофилы\", \"unit\": \"10^9/л\", \"value\": 0.03},\n    {\"name\": \"Моноциты\", \"unit\": \"10^9/л\", \"value\": 0.20},\n    {\"name\": \"Лимфоциты\", \"unit\": \"10^9/л\", \"value\": 1.91},\n    {\"name\": \"СОЭ (Вестергрен)\", \"unit\": \"мм/час\", \"value\": 2}\n  ]\n}"

	err = json.Unmarshal([]byte(content), &parsedContent)
	if err != nil {
		return entity.CheckupResults{}, err
	}

	return parsedContent.BuildEntity(), nil
}
