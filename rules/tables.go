//MIT License
//
//Copyright (c) 2023 arthur
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package rules

import (
	"github.com/arthur/banbrutes/config"
	"github.com/zngw/log"
)

func Check(text string, rule string) {
	err, ip, name, port := parse(text, rule)
	if err != nil {
		if err.Error() != "not tcp link" {
			log.Trace("net", err.Error())
		}

		return
	}

	rules(ip, name, port)

	return
}

func CheckFilter(text string, rule []config.RegFilter) {
	err, ip, name, port := parseFilter(text, rule)
	if err != nil {
		if err.Error() != "not tcp link" {
			log.Trace("net", err.Error())
		}

		return
	}

	rules(ip, name, port)

	return
}
