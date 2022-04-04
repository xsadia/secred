package internal

import "testing"

func TestNewWarehouseItemListFromCSV(t *testing.T) {
	t.Run("should return warehouse item list", func(t *testing.T) {
		validInput := [][]string{
			{"name", "max", "min", "quantity"},
			{"Rice", "50", "25", "27"},
		}

		wil, err := newWarehouseItemListFromCSV(validInput)

		if err != nil {
			t.Errorf("Expected no error. Got %q", err.Error())
		}

		if wil[0].Name != "rice" {
			t.Errorf("Expected first item's name to be \"rice\", got %q", wil[0].Name)
		}

		if wil[0].Max != 50 {
			t.Errorf("Expected first item's max to be 50, got %d", wil[0].Max)
		}

		if wil[0].Min != 25 {
			t.Errorf("Expected first item's min to be 25, got %d", wil[0].Min)
		}

		if wil[0].Quantity != 27 {
			t.Errorf("Expected first item's quantity to be 27, got %d", wil[0].Quantity)
		}

	})

	t.Run("Should return error if headers are not correct", func(t *testing.T) {
		invalidInput := [][]string{
			{"max", "min", "quantity", "babadook"},
			{"Rice", "50", "25", "27"},
		}

		_, err := newWarehouseItemListFromCSV(invalidInput)

		if err == nil {
			t.Error("Expected an error got none")
		}

		if err.Error() != badCSVHeadersError {
			t.Errorf("Expected error to be %q, got %q", badCSVHeadersError, err.Error())
		}
	})
}
