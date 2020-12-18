// Code generated by execgen; DO NOT EDIT.
// Copyright 2018 The Cockroach Authors.
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colexecagg

import (
	"strings"
	"unsafe"

	"github.com/cockroachdb/apd/v2"
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexec/execgen"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecbase/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
	"github.com/cockroachdb/errors"
)

func newSumOrderedAggAlloc(
	allocator *colmem.Allocator, t *types.T, allocSize int64,
) (aggregateFuncAlloc, error) {
	allocBase := aggAllocBase{allocator: allocator, allocSize: allocSize}
	switch t.Family() {
	case types.IntFamily:
		switch t.Width() {
		case 16:
			return &sumInt16OrderedAggAlloc{aggAllocBase: allocBase}, nil
		case 32:
			return &sumInt32OrderedAggAlloc{aggAllocBase: allocBase}, nil
		default:
			return &sumInt64OrderedAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.DecimalFamily:
		return &sumDecimalOrderedAggAlloc{aggAllocBase: allocBase}, nil
	case types.FloatFamily:
		return &sumFloat64OrderedAggAlloc{aggAllocBase: allocBase}, nil
	case types.IntervalFamily:
		return &sumIntervalOrderedAggAlloc{aggAllocBase: allocBase}, nil
	default:
		return nil, errors.Errorf("unsupported sum %s agg type %s", strings.ToLower(""), t.Name())
	}
}

type sumInt16OrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg apd.Decimal
	// col points to the output vector we are updating.
	col []apd.Decimal
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
	overloadHelper              execgen.OverloadHelper
}

var _ AggregateFunc = &sumInt16OrderedAgg{}

func (a *sumInt16OrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Decimal()
}

func (a *sumInt16OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	oldCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int16(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	newCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumInt16OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumInt16OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroDecimalValue
	a.foundNonNullForCurrentGroup = false
}

type sumInt16OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumInt16OrderedAgg
}

var _ aggregateFuncAlloc = &sumInt16OrderedAggAlloc{}

const sizeOfSumInt16OrderedAgg = int64(unsafe.Sizeof(sumInt16OrderedAgg{}))
const sumInt16OrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumInt16OrderedAgg{}))

func (a *sumInt16OrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumInt16OrderedAggSliceOverhead + sizeOfSumInt16OrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumInt16OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type sumInt32OrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg apd.Decimal
	// col points to the output vector we are updating.
	col []apd.Decimal
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
	overloadHelper              execgen.OverloadHelper
}

var _ AggregateFunc = &sumInt32OrderedAgg{}

func (a *sumInt32OrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Decimal()
}

func (a *sumInt32OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	oldCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int32(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	newCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumInt32OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumInt32OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroDecimalValue
	a.foundNonNullForCurrentGroup = false
}

type sumInt32OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumInt32OrderedAgg
}

var _ aggregateFuncAlloc = &sumInt32OrderedAggAlloc{}

const sizeOfSumInt32OrderedAgg = int64(unsafe.Sizeof(sumInt32OrderedAgg{}))
const sumInt32OrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumInt32OrderedAgg{}))

func (a *sumInt32OrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumInt32OrderedAggSliceOverhead + sizeOfSumInt32OrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumInt32OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type sumInt64OrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg apd.Decimal
	// col points to the output vector we are updating.
	col []apd.Decimal
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
	overloadHelper              execgen.OverloadHelper
}

var _ AggregateFunc = &sumInt64OrderedAgg{}

func (a *sumInt64OrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Decimal()
}

func (a *sumInt64OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	oldCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int64(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)

						{

							tmpDec := &_overloadHelper.TmpDec1
							tmpDec.SetInt64(int64(v))
							if _, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, tmpDec); err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	newCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumInt64OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumInt64OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroDecimalValue
	a.foundNonNullForCurrentGroup = false
}

type sumInt64OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumInt64OrderedAgg
}

var _ aggregateFuncAlloc = &sumInt64OrderedAggAlloc{}

const sizeOfSumInt64OrderedAgg = int64(unsafe.Sizeof(sumInt64OrderedAgg{}))
const sumInt64OrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumInt64OrderedAgg{}))

func (a *sumInt64OrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumInt64OrderedAggSliceOverhead + sizeOfSumInt64OrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumInt64OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type sumDecimalOrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg apd.Decimal
	// col points to the output vector we are updating.
	col []apd.Decimal
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &sumDecimalOrderedAgg{}

func (a *sumDecimalOrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Decimal()
}

func (a *sumDecimalOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	oldCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Decimal(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							_, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, &v)
							if err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							_, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, &v)
							if err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)

						{

							_, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, &v)
							if err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroDecimalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)

						{

							_, err := tree.ExactCtx.Add(&a.curAgg, &a.curAgg, &v)
							if err != nil {
								colexecerror.ExpectedError(err)
							}
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	newCurAggSize := tree.SizeOfDecimal(&a.curAgg)
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumDecimalOrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumDecimalOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroDecimalValue
	a.foundNonNullForCurrentGroup = false
}

type sumDecimalOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumDecimalOrderedAgg
}

var _ aggregateFuncAlloc = &sumDecimalOrderedAggAlloc{}

const sizeOfSumDecimalOrderedAgg = int64(unsafe.Sizeof(sumDecimalOrderedAgg{}))
const sumDecimalOrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumDecimalOrderedAgg{}))

func (a *sumDecimalOrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumDecimalOrderedAggSliceOverhead + sizeOfSumDecimalOrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumDecimalOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type sumFloat64OrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg float64
	// col points to the output vector we are updating.
	col []float64
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &sumFloat64OrderedAgg{}

func (a *sumFloat64OrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Float64()
}

func (a *sumFloat64OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Float64(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroFloat64Value

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							a.curAgg = float64(a.curAgg) + float64(v)
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroFloat64Value

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)

						{

							a.curAgg = float64(a.curAgg) + float64(v)
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroFloat64Value

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)

						{

							a.curAgg = float64(a.curAgg) + float64(v)
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroFloat64Value

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)

						{

							a.curAgg = float64(a.curAgg) + float64(v)
						}

						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumFloat64OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumFloat64OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroFloat64Value
	a.foundNonNullForCurrentGroup = false
}

type sumFloat64OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumFloat64OrderedAgg
}

var _ aggregateFuncAlloc = &sumFloat64OrderedAggAlloc{}

const sizeOfSumFloat64OrderedAgg = int64(unsafe.Sizeof(sumFloat64OrderedAgg{}))
const sumFloat64OrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumFloat64OrderedAgg{}))

func (a *sumFloat64OrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumFloat64OrderedAggSliceOverhead + sizeOfSumFloat64OrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumFloat64OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type sumIntervalOrderedAgg struct {
	orderedAggregateFuncBase
	// curAgg holds the running total, so we can index into the slice once per
	// group, instead of on each iteration.
	curAgg duration.Duration
	// col points to the output vector we are updating.
	col []duration.Duration
	// foundNonNullForCurrentGroup tracks if we have seen any non-null values
	// for the group that is currently being aggregated.
	foundNonNullForCurrentGroup bool
}

var _ AggregateFunc = &sumIntervalOrderedAgg{}

func (a *sumIntervalOrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Interval()
}

func (a *sumIntervalOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	var oldCurAggSize uintptr
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Interval(), vec.Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups and col to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		col := col
		if sel == nil {
			_ = groups[inputLen-1]
			_ = col.Get(inputLen - 1)
			if nulls.MaybeHasNulls() {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroIntervalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						//gcassert:bce
						v := col.Get(i)
						a.curAgg = a.curAgg.Add(v)
						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for i := 0; i < inputLen; i++ {

					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroIntervalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						//gcassert:bce
						v := col.Get(i)
						a.curAgg = a.curAgg.Add(v)
						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		} else {
			sel = sel[:inputLen]
			if nulls.MaybeHasNulls() {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroIntervalValue

							a.foundNonNullForCurrentGroup = false
						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = nulls.NullAt(i)
					if !isNull {
						v := col.Get(i)
						a.curAgg = a.curAgg.Add(v)
						a.foundNonNullForCurrentGroup = true
					}
				}
			} else {
				for _, i := range sel {

					if groups[i] {
						if !a.isFirstGroup {
							// If we encounter a new group, and we haven't found any non-nulls for the
							// current group, the output for this group should be null.
							if !a.foundNonNullForCurrentGroup {
								a.nulls.SetNull(a.curIdx)
							} else {
								a.col[a.curIdx] = a.curAgg
							}
							a.curIdx++
							a.curAgg = zeroIntervalValue

						}
						a.isFirstGroup = false
					}

					var isNull bool
					isNull = false
					if !isNull {
						v := col.Get(i)
						a.curAgg = a.curAgg.Add(v)
						a.foundNonNullForCurrentGroup = true
					}
				}
			}
		}
	},
	)
	var newCurAggSize uintptr
	if newCurAggSize != oldCurAggSize {
		a.allocator.AdjustMemoryUsage(int64(newCurAggSize - oldCurAggSize))
	}
}

func (a *sumIntervalOrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// null.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.col[outputIdx] = a.curAgg
	}
}

func (a *sumIntervalOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = zeroIntervalValue
	a.foundNonNullForCurrentGroup = false
}

type sumIntervalOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []sumIntervalOrderedAgg
}

var _ aggregateFuncAlloc = &sumIntervalOrderedAggAlloc{}

const sizeOfSumIntervalOrderedAgg = int64(unsafe.Sizeof(sumIntervalOrderedAgg{}))
const sumIntervalOrderedAggSliceOverhead = int64(unsafe.Sizeof([]sumIntervalOrderedAgg{}))

func (a *sumIntervalOrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(sumIntervalOrderedAggSliceOverhead + sizeOfSumIntervalOrderedAgg*a.allocSize)
		a.aggFuncs = make([]sumIntervalOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}
