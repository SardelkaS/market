package failure

type LogicError struct {
	code        int
	err         error
	description string
}

func (l LogicError) Code() int {
	return l.code
}

func (l LogicError) Description() string {
	return l.description
}

func (l LogicError) Error() string {
	return l.err.Error()
}

func (l LogicError) Wrap(err error) LogicError {
	l.err = err
	return l
}

var (
	// Basic
	ErrInput           = LogicError{code: 500, description: "bad request"}
	ErrNotFound        = LogicError{code: 404, description: "not found"}
	ErrAuth            = LogicError{code: 401, description: "unauthorized"}
	ErrGetUser         = LogicError{code: 500, description: "user not found"}
	ErrServiceNotFound = LogicError{code: 404, description: "service not found"}
	ErrForbidden       = LogicError{code: 403, description: "forbidden"}
	// Basket
	ErrAddProduct  = LogicError{code: 500, description: "error to add product to basket"}
	ErrIncrProduct = LogicError{code: 500, description: "error to increment product count"}
	ErrDecrProduct = LogicError{code: 500, description: "error to decrement product count"}
	ErrClearBasket = LogicError{code: 500, description: "error to clear basket"}
	ErrGetBasket   = LogicError{code: 500, description: "error to get basket"}
	// Feedback
	ErrSaveFeedback     = LogicError{code: 500, description: "error to save feedback"}
	ErrGetFeedback      = LogicError{code: 500, description: "error to get feedback"}
	ErrFeedbackNotFound = LogicError{code: 500, description: "feedback not found"}
	ErrRemoveFeedback   = LogicError{code: 500, description: "error to remove feedback"}
	ErrYourFeedback     = LogicError{code: 500, description: "it's your feedback"}
	ErrLikeFeedback     = LogicError{code: 500, description: "error to like feedback"}
	ErrUnlikeFeedback   = LogicError{code: 500, description: "error to unlike feedback"}
	ErrFeedbackLiked    = LogicError{code: 500, description: "feedback is already liked"}
	// Order
	ErrCreateOrder              = LogicError{code: 500, description: "error to create order"}
	ErrGetOrder                 = LogicError{code: 500, description: "error to get order"}
	ErrAttachProduct            = LogicError{code: 500, description: "error to attach product"}
	ErrDetachProduct            = LogicError{code: 500, description: "error to detach product"}
	ErrUpdateOrderProductsCount = LogicError{code: 500, description: "error to update products count"}
	ErrUpdateOrderStatus        = LogicError{code: 500, description: "error to update order status"}
	ErrGetOrderProducts         = LogicError{code: 500, description: "error to get order products"}
	// Categories
	ErrGetCategories = LogicError{code: 500, description: "error to get categories"}
	// Manufacturers
	ErrGetManufacturers = LogicError{code: 500, description: "error to get manufacturers"}
	// Sexes
	ErrGetSexes = LogicError{code: 500, description: "error to get sexes"}
	// Countries
	ErrGetCountries = LogicError{code: 500, description: "error to get countries"}
	// Product
	ErrGetProduct                = LogicError{code: 500, description: "error to get product"}
	ErrFindProduct               = LogicError{code: 500, description: "error to find product"}
	ErrLikeProduct               = LogicError{code: 500, description: "error to like product"}
	ErrUnlikeProduct             = LogicError{code: 500, description: "error to unlike product"}
	ErrCheckProductLiked         = LogicError{code: 500, description: "error to check is product liked"}
	ErrGetProductStars           = LogicError{code: 500, description: "error to get product stars"}
	ErrUpdateProductsCount       = LogicError{code: 500, description: "error to update products count"}
	ErrViewProduct               = LogicError{code: 500, description: "error to view product"}
	ErrGetRecentlyViewedProducts = LogicError{code: 500, description: "error to get recently viewed products"}
	ErrGetBoughtProducts         = LogicError{code: 500, description: "error to get bought products"}
)
