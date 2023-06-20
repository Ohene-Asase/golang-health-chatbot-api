package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/reiver/go-porterstemmer"

)

type IntentData struct {
	Questions []string `json: "questions"`
	Answers   []string `json:answers"`
}

var intents = IntentData{
	Questions: []string{
		"Mouth Bitter",
		"Headache",
		"Cough and Sore Throat",
		"Fever and Body Aches",
		"Shortness of Breath",
		"Abdominal Pain and Diarrhea",
		"Joint Pain and Swelling",
		"Fatigue and Weakness",
		"Hello",
	},
	Answers: []string{
		"You might have Malaria",
		"It could be due to stress or a tension headache",
		"You may be suffering from a common cold or flu",
		"It could be a sign of influenza or dengue fever",
		"It could indicate a respiratory infection or asthma",
		"You may have food poisoning or a stomach virus",
		"It could be a symptom of arthritis or rheumatoid arthritis",
		"It may be due to lack of sleep or anemia",
		"Hi, How may I help you",
	},
}

func main() {
	manager := NewNlpManager()
	err := manager.Train()
	if err != nil {
		log.Fatal("Error occured while training the model:", err)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world")
	})

	router.GET("/api/chatbot/:message", func(c *gin.Context) {
		message := c.Param("message")
		answer := manager.Process(message)
		c.JSON(http.StatusOK, gin.H{"answer": answer})
	})

	err = router.Run(":3000")
	if err != nil {
		log.Fatal("Error occured while starting the server:", err)
	}

}

type NlpManager struct {
	Languages     []string
	IntentManager map[string]string
}

func NewNlpManager() *NlpManager {
	manager := &NlpManager{
		Languages:     []string{"en"},
		IntentManager: make(map[string]string),
	}

	for i, question := range intents.Questions {
		intent := fmt.Sprintf("intent_%d", i+1)
		manager.addDocument(question, intent)
		manager.addAnswer(intent, intents.Answers[i])
	}

	return manager
}

func (m *NlpManager) addDocument(question, intent string) {
	m.IntentManager[normalizeText(question)] = intent
}

func (m *NlpManager) addAnswer(intent, answer string) {
	m.IntentManager[intent] = answer
}

func (m *NlpManager) Train() error {
	fmt.Println("Model trained successfully")
	return nil
}

func (m *NlpManager) Process(message string) string {
	intent, ok := m.IntentManager[normalizeText(message)]
	if !ok || !strings.HasPrefix(intent, "intent_") {
		return "Sorry, I don't"
	}

	index, _ := strconv.Atoi(strings.TrimPrefix(intent, "intent_"))
	if index >= 1 && index <= len(intents.Answers) {
		return intents.Answers[index-1]
	}
	return ""
}

func normalizeText(text string) string {
	words := strings.Fields(text)
	normalizedWords := make([]string, len(words))

	for i, word := range words {
		stemmedWord := porterstemmer.StemString(word)
		normalizedWords[i] = stemmedWord
	}

	normalizedText := strings.Join(normalizedWords, " ")
	return normalizedText
}
