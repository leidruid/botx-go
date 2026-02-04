package bot

import (
	"context"

	"github.com/google/uuid"
	"github.com/leidruid/botx-go/models"
	"github.com/leidruid/botx-go/optional"
)

// SearchUserByHUID returns user profile by HUID.
func (b *Bot) SearchUserByHUID(ctx context.Context, botID uuid.UUID, huid uuid.UUID) (models.UserFromSearch, error) {
	return b.Client.SearchUserByHUID(ctx, botID, huid)
}

// SearchUserByEmail returns user profile by email.
func (b *Bot) SearchUserByEmail(ctx context.Context, botID uuid.UUID, email string) (models.UserFromSearch, error) {
	return b.Client.SearchUserByEmail(ctx, botID, email)
}

// SearchUserByEmails returns user profiles by emails.
func (b *Bot) SearchUserByEmails(ctx context.Context, botID uuid.UUID, emails []string) ([]models.UserFromSearch, error) {
	return b.Client.SearchUserByEmails(ctx, botID, emails)
}

// SearchUserByLogin returns user profile by AD login/domain.
func (b *Bot) SearchUserByLogin(ctx context.Context, botID uuid.UUID, adLogin string, adDomain string) (models.UserFromSearch, error) {
	return b.Client.SearchUserByLogin(ctx, botID, adLogin, adDomain)
}

// SearchUserByOtherID returns user profile by other id.
func (b *Bot) SearchUserByOtherID(ctx context.Context, botID uuid.UUID, otherID string) (models.UserFromSearch, error) {
	return b.Client.SearchUserByOtherID(ctx, botID, otherID)
}

// UpdateUserProfile updates user profile fields.
func (b *Bot) UpdateUserProfile(
	ctx context.Context,
	botID uuid.UUID,
	userHUID uuid.UUID,
	avatar optional.Optional[models.OutgoingAttachment],
	name optional.Optional[string],
	publicName optional.Optional[string],
	company optional.Optional[string],
	companyPosition optional.Optional[string],
	description optional.Optional[string],
	department optional.Optional[string],
	office optional.Optional[string],
	manager optional.Optional[string],
) error {
	return b.Client.UpdateUserProfile(ctx, botID, userHUID, avatar, name, publicName, company, companyPosition, description, department, office, manager)
}

// UsersAsCSV downloads users list as CSV and parses it.
func (b *Bot) UsersAsCSV(ctx context.Context, botID uuid.UUID, ctsUser bool, unregistered bool, botx bool) ([]models.UserFromCSV, error) {
	return b.Client.UsersAsCSV(ctx, botID, ctsUser, unregistered, botx)
}
