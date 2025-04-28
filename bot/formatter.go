package bot

import "fmt"

func FormatResponse(intent string, data map[string]interface{}) string {
	switch intent {
	case "ROAS":
		roas, ok := data["roas"].(float64)
		if !ok {
			return "Sorry, ROAS data is not available for your campaign."
		}
		return fmt.Sprintf("The ROAS for your campaign is %.2f.", roas)

	case "CTR":
		ctr, ok := data["ctr"].(float64)
		if !ok {
			return "Sorry, CTR data is not available for your campaign."
		}
		return fmt.Sprintf("The CTR for your campaign is %.2f%%.", ctr*100)

	case "Spend":
		spend, ok := data["spend"].(float64)
		if !ok {
			return "Sorry, Spend data is not available."
		}
		return fmt.Sprintf("You spent $%.2f on your campaigns.", spend)

	default:
		return "Sorry, I didn't understand your request."
	}
}
