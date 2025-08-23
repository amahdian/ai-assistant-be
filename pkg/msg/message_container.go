package msg

import (
	"fmt"
	"strings"
)

type MessageContainer struct {
	messages map[string][]*Message
	hasError bool
}

func NewMessageContainer() *MessageContainer {
	return &MessageContainer{
		messages: make(map[string][]*Message),
	}
}

func (mc *MessageContainer) Union(otherMessageContainers ...*MessageContainer) {
	for _, otherMessageContainer := range otherMessageContainers {
		for group, messages := range otherMessageContainer.messages {
			for _, message := range messages {
				switch message.Level {
				case Fatal:
					mc.AddFatal(group, message.Text)
				case Error:
					mc.AddError(group, message.Text)
				case Warning:
					mc.AddWarning(group, message.Text)
				case Info:
					mc.AddInfo(group, message.Text)
				}
			}
		}
	}
}

func (mc *MessageContainer) Count() int {
	count := 0
	for _, messages := range mc.messages {
		count = count + len(messages)
	}

	return count
}

func (mc *MessageContainer) ErrorCount() int {
	count := 0
	for _, messages := range mc.messages {
		for _, message := range messages {
			if message.Level == Error {
				count++
			}
		}
	}

	return count
}

func (mc *MessageContainer) HasError() bool {
	return mc.hasError
}

func (mc *MessageContainer) GetAll() map[string][]*Message {
	return mc.messages
}

func (mc *MessageContainer) getMessagesByLevel(level MessageLevel) []*Message {
	var messageArray []*Message
	for _, messages := range mc.messages {
		for _, message := range messages {
			if message.Level == level {
				messageArray = append(messageArray, message)
			}
		}
	}

	return messageArray
}

func (mc *MessageContainer) GetErrors() []*Message {
	return mc.getMessagesByLevel(Error)
}

func (mc *MessageContainer) GetWarnings() []*Message {
	return mc.getMessagesByLevel(Warning)
}

func (mc *MessageContainer) GetInfos() []*Message {
	return mc.getMessagesByLevel(Info)
}

func (mc *MessageContainer) AddFatal(group string, messageText string) {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  messageText,
		Level: Fatal,
	})
	mc.hasError = true
}

func (mc *MessageContainer) AddError(group string, messageText string) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  messageText,
		Level: Error,
	})
	mc.hasError = true
	return mc
}

func (mc *MessageContainer) AddWarning(group string, messageText string) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  messageText,
		Level: Warning,
	})
	return mc
}

func (mc *MessageContainer) AddInfo(group string, messageText string) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  messageText,
		Level: Info,
	})
	return mc
}

func (mc *MessageContainer) AddErrorf(group string, format string, parameters ...interface{}) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  fmt.Sprintf(format, parameters...),
		Level: Error,
	})
	mc.hasError = true
	return mc
}

func (mc *MessageContainer) AddWarningf(group string, format string, parameters ...interface{}) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  fmt.Sprintf(format, parameters...),
		Level: Warning,
	})
	return mc
}

func (mc *MessageContainer) AddInfof(group string, format string, parameters ...interface{}) *MessageContainer {
	mc.messages[group] = append(mc.messages[group], &Message{
		Text:  fmt.Sprintf(format, parameters...),
		Level: Info,
	})
	return mc
}

func MakePlainText(message string) string {
	newMessage := strings.ReplaceAll(message, string('"'), "`")
	newMessage = strings.ReplaceAll(message, string("'"), "`")
	newMessage = strings.ReplaceAll(newMessage, "\r\n", " ")
	newMessage = strings.ReplaceAll(newMessage, "\n", " ")
	newMessage = strings.ReplaceAll(newMessage, "\r", " ")
	newMessage = strings.ReplaceAll(newMessage, "\t", " ")
	return newMessage
}
