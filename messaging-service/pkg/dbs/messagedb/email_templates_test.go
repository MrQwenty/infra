package messagedb

import (
	"testing"

	"github.com/influenzanet/messaging-service/pkg/types"
)

func TestEmailTemplatesDB(t *testing.T) {
	t1 := types.EmailTemplate{
		MessageType:     "test-type",
		StudyKey:        "test-study",
		DefaultLanguage: "en",
	}
	t.Run("save not existing template", func(t *testing.T) {
		_, err := testDBService.SaveEmailTemplate(testInstanceID, t1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("save existing template", func(t *testing.T) {
		t1.DefaultLanguage = "de"
		res, err := testDBService.SaveEmailTemplate(testInstanceID, t1)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if res.DefaultLanguage != "de" {
			t.Errorf("unexpected result: %v", res)
		}
	})

	t.Run("delete existing template", func(t *testing.T) {
		err := testDBService.DeleteEmailTemplate(testInstanceID, t1.MessageType, t1.StudyKey)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})

	t.Run("delete not existing template", func(t *testing.T) {
		err := testDBService.DeleteEmailTemplate(testInstanceID, t1.MessageType, t1.StudyKey)
		if err == nil {
			t.Error("should fail when template already deleted")
			return
		}
	})

	testTemplates := []types.EmailTemplate{
		{
			MessageType:     "t1",
			StudyKey:        "test-study",
			DefaultLanguage: "en",
		},
		{
			MessageType:     "t2",
			DefaultLanguage: "de",
		},
		{
			MessageType:     "t2",
			StudyKey:        "for study",
			DefaultLanguage: "fr",
		},
	}
	for _, temp := range testTemplates {
		_, err := testDBService.SaveEmailTemplate(testInstanceID, temp)
		if err != nil {
			t.Errorf("unexpected error when creating test templates: %v", err)
			return
		}
	}

	t.Run("find template by message type", func(t *testing.T) {
		res, err := testDBService.FindEmailTemplateByType(testInstanceID, "t2", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if res.DefaultLanguage != "de" || res.StudyKey != "" {
			t.Errorf("unexpected template found: %v", res)
		}
	})

	t.Run("find template by message type and study key", func(t *testing.T) {
		res, err := testDBService.FindEmailTemplateByType(testInstanceID, "t2", "for study")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if res.DefaultLanguage != "fr" || res.StudyKey != "for study" {
			t.Errorf("unexpected template found: %v", res)
		}
	})

	t.Run("find all templates", func(t *testing.T) {
		res, err := testDBService.FindAllEmailTempates(testInstanceID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		if len(res) != 3 {
			t.Errorf("unexpected number of templates found: %d", len(res))
		}
	})

}
