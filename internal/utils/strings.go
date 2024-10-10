package utils

import (
	"strconv"
	"strings"
)

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

func levenshtein(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	// Create a 2D slice to store the distances
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize the first row and column
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Calculate Levenshtein distance
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+1)
			}
		}
	}

	return dp[m][n]
}

var romanSymbols = []struct {
	Value  int16
	Symbol string
}{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func Int2Roman(num int16) string {
	roman := ""
	for _, symbol := range romanSymbols {
		for num >= symbol.Value {
			roman += symbol.Symbol
			num -= symbol.Value
		}
	}
	return roman
}

func IntSlice2String(slice []int) string {
	str := ""
	for _, i := range slice {
		str += strconv.Itoa(i) + ","
	}
	return strings.TrimSuffix(str, ",")
}
