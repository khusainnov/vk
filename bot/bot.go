package bot

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	errClosedCh    = errors.New("closed update channel")
	errUnsupported = errors.New("unsupported command type")
)

var (
	whatForMsg    = "Я создан в качестве тестового задания в отборе на одну из стажировок о которых собираюсь рассказать. Если тебе интересны такие события, то оставайся со мной и смотри что еще я умею."
	internshipMsg = "Оплачиваемая стажировка в VK — это твой шанс попасть в компанию. Наравне со всеми сотрудниками ты будешь работать над продуктами, которыми пользуются миллионы людей.\n\nЭто отличное решение, если ты учишься в вузе, у тебя еще нет опыта и ты пока не знаешь, чем хочешь заниматься в будущем.  В течение 2–5 месяцев ты поработаешь над реальными кейсами под руководством наставников и поймёшь, что подходит именно тебе\n\n"
)

const (
	startCmd      = "/start"
	menuCmd       = "/menu"
	whatForImCmd  = "/whatfor"
	internshipCmd = "/internship"
	routesCmd     = "/routes"
)

func parseCommand(msg string) string {
	var rsp bool

	if rsp = strings.Contains(msg, startCmd); rsp {
		return startCmd
	}
	if rsp = strings.Contains(msg, menuCmd); rsp {
		return menuCmd
	}
	if rsp = strings.Contains(msg, whatForImCmd); rsp {
		return whatForImCmd
	}
	if rsp = strings.Contains(msg, internshipCmd); rsp {
		return internshipCmd
	}
	if rsp = strings.Contains(msg, routesCmd); rsp {
		return routesCmd
	}

	return ""
}

func readFile(name string) ([]byte, error) {
	data, err := os.ReadFile(fmt.Sprintf("./third_party/assets/%s", name))
	if err != nil {
		return nil, fmt.Errorf("cannot read image, %w", err)
	}

	return data, nil
}
