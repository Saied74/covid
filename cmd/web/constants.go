package main

// //States is exported for templates
// var States = StatesType{
// 	State: []string{"Alabama", "Alaska", "Arizona", "Arkansas", "California",
// 		"Colorado", "Connecticut", "Delaware", "Florida", "Georgia",
// 		"Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas",
// 		"Kentucky", "Louisiana", "Maine", "Maryland", "Massachusetts",
// 		"Michigan", "Minnesota", "Mississippi", "Missouri", "Montana",
// 		"Nebraska", "Nevada", "New Hampshire", "New Jersey",
// 		"New Mexico", "New York", "North Carolina", "North Dakota",
// 		"Ohio", "Oklahoma", "Oregon", "Pennsylvania", "Rhode Island",
// 		"South Carolina", "South Dakota", "Tennessee", "Texas",
// 		"Utah", "Vermont", "Virginia", "Washington", "West Virginia",
// 		"Wisconsin", "Wyoming",
// 	},
// 	Short: []string{"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
// 		"HI", "ID", "IL", "IM", "IA", "KS", "KY", "LA", "ME", "MD",
// 		"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
// 		"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
// 		"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
// 	},
// 	Fields: []string{"positive", "negative", "pending", "hospitalized",
// 		"death", "total", "hash", "dateChecked", "totalTestResults",
// 		"flips", "deathIncrease", "hospitalizedIncrease",
// 		"positiveIncrease", "totalTestResultsIncrease",
// 	},
// 	Xdata: []string{"1", "2", "3", "4", "5"},
// 	// Ydata: []string{"1", "2", "4", "8", "16"},
// }

var states = []string{"Alabama", "Alaska", "Arizona", "Arkansas", "California",
	"Colorado", "Connecticut", "Delaware", "Florida", "Georgia",
	"Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas",
	"Kentucky", "Louisiana", "Maine", "Maryland", "Massachusetts",
	"Michigan", "Minnesota", "Mississippi", "Missouri", "Montana",
	"Nebraska", "Nevada", "New Hampshire", "New Jersey",
	"New Mexico", "New York", "North Carolina", "North Dakota",
	"Ohio", "Oklahoma", "Oregon", "Pennsylvania", "Rhode Island",
	"South Carolina", "South Dakota", "Tennessee", "Texas",
	"Utah", "Vermont", "Virginia", "Washington", "West Virginia",
	"Wisconsin", "Wyoming",
}

var short = []string{"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
	"HI", "ID", "IL", "IM", "IA", "KS", "KY", "LA", "ME", "MD",
	"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
	"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
	"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
}

var fields = []string{"positive", "negative", "pending", "hospitalized",
	"death", "total", "hash", "dateChecked", "totalTestResults",
	"flips", "deathIncrease", "hospitalizedIncrease",
	"positiveIncrease", "totalTestResultsIncrease",
}
