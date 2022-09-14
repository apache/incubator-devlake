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

// Convert converts a raw error to an Error. Similar to Type.WrapRaw with a type of Default, but will not wrap the error if it's already of type Error.
// Passing in nil will return nil.
func Convert(err error) Error {
	return Default.wrapRaw(err, false, withStackOffset(1))
}

// Convert01 an extension of Convert that allows passing in one extra arg. Useful for inlining the conversion of multivalued returns.
func Convert01[T1 any](t1 T1, err error) (T1, Error) {
	return t1, Default.wrapRaw(err, false, withStackOffset(1))
}

// Convert001 like Convert01, but with two extra args
func Convert001[T1 any, T2 any](t1 T1, t2 T2, err error) (T1, T2, Error) {
	return t1, t2, Default.wrapRaw(err, false, withStackOffset(1))
}

// Convert0001 like Convert01, but with three extra args
func Convert0001[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, err error) (T1, T2, T3, Error) {
	return t1, t2, t3, Default.wrapRaw(err, false, withStackOffset(1))
}

// Convert00001 like Convert01, but with four extra args
func Convert00001[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) (T1, T2, T3, T4, Error) {
	return t1, t2, t3, t4, Default.wrapRaw(err, false, withStackOffset(1))
}
