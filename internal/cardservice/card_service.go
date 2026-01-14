package cardservice
import "transaction-manager/internal/domain"
type CardEventService interface {
	HandleCardAuth(event domain.CardEvent) error
	// HandleCardSettlement(event domain.CardEvent) error
	// HandleCardRelease(event domain.CardEvent) error
}
