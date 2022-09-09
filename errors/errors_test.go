/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestCrdbErrorImpl(t *testing.T) {
	err := f1()
	lakeErr := AsLakeErrorType(err)
	require.NotNil(t, lakeErr)
	t.Run("full_error", func(t *testing.T) {
		fmt.Printf("======================Full Error=======================: \n%v\n\n\n", err)
		require.Equal(t, err.Error(), lakeErr.Error())
	})
	t.Run("raw_message", func(t *testing.T) {
		msg := lakeErr.Messages().Format(RawMessageType)
		require.NotEqual(t, err.Error(), msg)
		fmt.Printf("======================Raw Message=======================: \n%s\n\n\n", msg)
		msgParts := strings.Split(msg, "\ncaused by: ")
		expectedParts := []string{
			"f1 error (404)",
			"f2 error (404)",
			"f3 error (400)",
			os.ErrNotExist.Error() + " (400)",
		}
		require.Equal(t, expectedParts, msgParts)
	})
	t.Run("type_conversion", func(t *testing.T) {
		e := lakeErr.As(NotFound)
		require.Equal(t, NotFound, e.GetType())
		e = lakeErr.As(BadInput)
		require.Equal(t, NotFound, e.GetType())
		e = lakeErr.As(Internal)
		require.Nil(t, e)
	})
	t.Run("type_casting", func(t *testing.T) {
		require.True(t, errors.Is(lakeErr, os.ErrNotExist))
	})
	t.Run("combine_errors_type", func(t *testing.T) {
		err = Unauthorized.Combine([]error{err, err})
		lakeErr = AsLakeErrorType(err)
		require.NotNil(t, lakeErr)
		e := lakeErr.As(Unauthorized)
		require.Equal(t, Unauthorized, e.GetType())
		e = lakeErr.As(NotFound)
		require.Nil(t, e)
		e = lakeErr.As(BadInput)
		require.Nil(t, e)
		require.False(t, errors.Is(lakeErr, os.ErrNotExist))
	})
	t.Run("error convert", func(t *testing.T) {
		rawErr := errors.New("test error")
		err := Convert(rawErr)
		require.Equal(t, rawErr, err.Unwrap())
		require.Equal(t, Internal.GetHttpCode(), err.GetType().GetHttpCode())
		baseErr := BadInput.Wrap(rawErr, "wrapped")
		err2 := Convert(baseErr)
		require.Same(t, baseErr, err2)
		require.Equal(t, "wrapped (400)", err2.Messages().Get(RawMessageType))
		require.Same(t, rawErr, err2.Unwrap())
		err3 := Default.WrapRaw(baseErr)
		require.NotSame(t, baseErr, err3)
		require.Equal(t, "wrapped (400)", err3.Messages().Get(RawMessageType))
		require.Equal(t, "wrapped (400)", err3.Messages().Get(UserMessageType))
		require.Same(t, baseErr, err3.Unwrap())
	})
}

func f1() Error {
	err := f2()
	return Default.Wrap(err, "f1 error")
}

func f2() Error {
	err := f3()
	return NotFound.Wrap(err, "f2 error")
}

func f3() Error {
	err := f4()
	return Default.Wrap(err, "f3 error")
}

func f4() Error {
	err := f5()
	return BadInput.WrapRaw(err)
}

func f5() error {
	return os.ErrNotExist
}
