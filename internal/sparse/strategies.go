package sparse

import (
	"math/bits"
)

/**
 *  Word and block <b>and</b> strategy.
 */
var andStrategy = new(andStrategyType)

/**
 *  Word and block <b>andNot</b> strategy.
 */
var andNotStrategy = new(andNotStrategyType)

/**
 *  Word and block <b>clear</b> strategy.
 */
var clearStrategy = new(clearStrategyType)

/**
 *  Word and block <b>copy</b> strategy.
 */
var copyStrategy = new(copyStrategyType)

/**
 *  Word and block <b>flip</b> strategy.
 */
var flipStrategy = new(flipStrategyType)

/**
 *  Word and block <b>or</b> strategy.
 */
var orStrategy = new(orStrategyType)

/**
 *  Word and block <b>set</b> strategy.
 */
var setStrategy = new(setStrategyType)

/**
 *  Word and block <b>xor</b> strategy.
 */
var xorStrategy = new(xorStrategyType)

/**
 *  And of two sets. Where the <i>a</i> set is zero, it remains zero (i.e.,
 *  without entries or with zero words). Similarly, where the <i>b</i> set is
 *  zero, the <i>a</i> becomes zero (i.e., without entries).
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched virtual "entries" of the <i>b</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these are all false, hence
 *  the result of the "and" operation will be to make all these entries in this
 *  set to become false--hence just remove them, and then scan only those
 *  entries that could match entries in the bit set<i>b</i>. This clearing of
 *  the remainder of the <i>a</i> set is accomplished by selecting both
 *  <i>F_OP_X_EQ_F</i> and <i>X_OP_F_EQ_F</i>.
 *
 *  <pre>
 *  and| 0 1
 *    0| 0 0
 *    1| 0 1 <pre>
 */
type andStrategyType struct{}

func (st andStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cValueOpFalseEqFalse + cValueOpFalseEqFalse
}

func (st andStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	return true
}

func (st andStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3]&b3[u3] | ^mask
	return a3[u3] == 0
}

func (st andStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] & b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}

func (st andStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  AndNot of two sets. Where the <i>a</i> set is zero, it remains zero
 *  (i.e., without entries or with zero words). On the other hand, where the
 *  <i>b</i> set is zero, the <i>a</i> remains unchanged.
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched virtual "entries" of the <i>b</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these are all false, hence
 *  the result of the "and" operation will be to make all these entries in this
 *  set to become false--hence just remove them, and then scan only those
 *  entries that could match entries in the bit set<i>b</i>. This clearing of
 *  the remainder of the <i>a</i> set is accomplished by selecting both
 *  <i>F_OP_X_EQ_F</i> and <i>X_OP_F_EQ_F</i>.
 *
 *  <pre>
 * andNot| 0 1
 *      0| 0 0
 *      1| 1 0 <pre>
 */
type andNotStrategyType struct{}

func (st andNotStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cFalseOpValueEqFalse + cValueOpFalseEqValue
}

func (st andNotStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	return true
}

func (st andNotStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] & ^(b3[u3] & mask)
	return a3[u3] == 0
}

func (st andNotStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] & ^b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st andNotStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Clear clears bits in the <i>a</i> set.
 *
 * <pre>
 * clear| 0 1
 *     0| 0 0
 *     1| 0 0 <pre>
 */
type clearStrategyType struct{}

func (st clearStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cFalseOpValueEqFalse
}

func (st clearStrategyType) start(b *BitSet) bool {
	return true
}

func (st clearStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] & ^mask
	return a3[u3] == 0
}

func (st clearStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	if u3 != 0 || v3 != cLength3 {
		for w3 := u3; w3 != v3; w3 = w3 + 1 {
			a3[w3] = 0
		}
	}
	return true
}
func (st clearStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Copies the needed parts of the <i>b</i> set to the <i>a</i> set.
 *
 * <pre>
 * get| 0 1
 *   0| 0 1
 *   1| 0 1 <pre>
 */
type copyStrategyType struct{}

func (st copyStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cValueOpFalseEqFalse
}

func (st copyStrategyType) start(b *BitSet) bool {
	return true
}

func (st copyStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = b3[u3] & mask
	return a3[u3] == 0
}

func (st copyStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st copyStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Equals compares bits in the <i>a</i> set with those in the <i>b</i> set.
 *  None of the values in either set are changed, although the <i>a</i> set
 *  may have all zero level 3 blocks replaced by null references (and
 *  similarly at level 2).
 *
 * <pre>
 * equals| 0 1
 *      0| 0 -
 *      1| - - <pre>
 */
type equalsStrategyType struct {
	result bool
}

func (st equalsStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse
}

func (st *equalsStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	st.result = true
	return false
	/*  Equals does not change the content of the set, hence hash need
	    not be reset. */
}

func (st *equalsStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	st.result = st.result && ((word & mask) == (b3[u3] & mask))
	return word == 0
}

func (st *equalsStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		word := a3[w3]
		st.result = st.result && word == b3[w3]
		isZero = isZero && word == 0
	}
	return
}

func (st equalsStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Flip inverts the bits of the <i>a</i> set within the given range.
 *
 * <pre>
 * flip| 0 1
 *    0| 1 1
 *    1| 0 0 <pre>
 */
type flipStrategyType struct{}

func (st flipStrategyType) properties() int32 {
	return 0
}

func (st flipStrategyType) start(b *BitSet) bool {
	return true
}

func (st flipStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] ^ mask
	return a3[u3] == 0
}

func (st flipStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ wordType(^wordType(0))
		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st flipStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Intersect has a true result if any word in the <i>a</i> set has a bit
 *  in common with the <i>b</i> set. During the scan of the <i>a</i> set
 *  blocks (and areas) that are all zero may be replaced with empty blocks
 *  and areas (null references), but the value of the set is not changed
 *  (which is why X_OP_F_EQ_F is not selected, since this would cause
 *  parts of the <i>a</i> set to be zero-ed out).
 *
 * <pre>
 * intersect| 0 1
 *         0| 0 0
 *         1| 1 1 <pre>
 */
type intersectsStrategyType struct {
	result bool
}

func (st intersectsStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cFalseOpValueEqFalse
}

func (st *intersectsStrategyType) start(b *BitSet) bool {
	if b == nil {
		panic("b is nil")
	}
	st.result = false
	return false
	/*  Intersect does not change the content of the set, hence hash need
	    not be reset. */
}

func (st *intersectsStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	st.result = st.result || ((word & b3[u3] & mask) != 0)
	return word == 0
}

func (st *intersectsStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		word := a3[w3]
		st.result = st.result || ((word & b3[w3]) != 0)
		isZero = isZero && word == 0
	}
	return
}
func (st intersectsStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

/**
 *  Or of two sets. Where the <i>a</i> set is one, it remains one. Similarly,
 *  where the <i>b</i> set is one, the <i>a</i> becomes one. If both sets have
 *  zeros in corresponding places, a zero results. Whole blocks or areas that
 *  are or become zero are replaced by null arrays.
 *  <p>
 *  If level1 of the <i>a</i> set is longer than level1 of the bit set
 *  <i>b</i>, then the unmatched entries of the <i>a</i> set (beyond
 *  the actual length of <i>b</i>) corresponding to these remain unchanged. *
 *  <pre>
 *   or| 0 1
 *    0| 0 1
 *    1| 1 1 <pre>
 */

type orStrategyType struct{}

func (st orStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cValueOpFalseEqValue
}

func (st orStrategyType) start(b *BitSet) bool {
	return true
}

func (st orStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] | (b3[u3] & mask)
	return a3[u3] == 0
}

func (st orStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] | b3[w3]

		isZero = isZero && a3[w3] == 0
	}
	return
}
func (st orStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Set creates entries everywhere within the range. Hence no empty level2
 *  areas or level3 blocks are ignored, and no empty (all zero) blocks are
 *  returned.
 *
 *  <pre>
 * set| 0 1
 *   0| 1 1
 *   1| 1 1 <pre>
 */

type setStrategyType struct{}

func (st setStrategyType) properties() int32 {
	return 0
}

func (st setStrategyType) start(b *BitSet) bool {
	return true
}

func (st setStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] | mask
	return a3[u3] == 0
}

func (st setStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = ^wordType(0)
	}
	isZero = false
	return
}
func (st setStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}

//-----------------------------------------------------------------------------
/**
 *  Update the seven statistics that are computed for each set. These are
 *  updated by calling <i>statisticsUpdate</i>, which uses this strategy.
 *
 *  <pre>
 *  update| 0 1
 *       0| 0 0
 *       1| 1 1 <pre>
 *
 * @see SparseBitSet#statisticsUpdate()
 */
type updateStrategyType struct {
	/**
	 *  Working space for find the size and length of the bit set. Holds copy of
	 *  the first non-empty word in the set.
	 */
	wordMin wordType
	/**
	 *  Working space for find the size and length of the bit set. Holds a copy
	 *  of the last non-empty word in the set.
	 */
	wordMax wordType

	/**
	 *  Working space for find the hash value of the bit set. Holds the
	 *  current state of the computation of the hash value. This value is
	 *  ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	hash uint64
	/**
	 *  Working space for find the size and length of the bit set. Holds the
	 *  index of the first non-empty word in the set.
	 */
	wMin int32

	/**
	 *  Working space for find the size and length of the bit set. Holds the
	 *  index of the last non-empty word in the set.
	 */
	wMax int32

	/**
	 *  Working space for keeping count of the number of non-zero words in the
	 *  bit set. Holds the current state of the computation of the count. This
	 *  value is ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	count int32

	/**
	 *  Working space for counting the number of non-zero bits in the bit set.
	 *  Holds the current state of the computation of the cardinality.This
	 *  value is ultimately transferred to the Cache object.
	 *
	 * @see SparseBitSet.Cache
	 */
	cardinality int32
}

func (st updateStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cFalseOpValueEqFalse
}

func (st *updateStrategyType) start(b *BitSet) bool {
	st.hash = 1234     // Magic number
	st.wMin = -1       // index of first non-zero word
	st.wordMin = 0     // word at that index
	st.wMax = 0        // index of last non-zero word
	st.wordMax = 0     // word at that index
	st.count = 0       // count of non-zero words in whole set
	st.cardinality = 0 // count of non-zero bits in the whole set
	return false
}

func (st *updateStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	word := a3[u3]
	word1 := word & mask
	if word1 != 0 {
		st.compute(base+u3, word1)
	}
	return word == 0
}

func (st *updateStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	isZero = true //  Presumption
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		if word := a3[w3]; word != 0 {
			isZero = false
			st.compute(base+w3, word)
		}
	}
	return isZero
}

func (st *updateStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {
	cache.a2Count = a2Count
	cache.a3Count = a3Count
	cache.count = st.count
	cache.cardinality = st.cardinality
	cache.length = (st.wMax+1)*cLength4 - int32(bits.LeadingZeros(uint(st.wordMax)))
	cache.size = cache.length - st.wMin*cLength4 - int32(bits.LeadingZeros(uint(st.wordMin)))
	cache.hash = ((st.hash >> cIntegerSize) ^ st.hash)
}

func (st *updateStrategyType) compute(index int32, word wordType) {
	/*  Count the number of actual words being used. */
	st.count++
	/*  Continue to accumulate the hash value of the set. */
	st.hash = st.hash ^ (uint64(word) * uint64(index+1))
	/*  The first non-zero word contains the first actual bit of the
	    set. The location of this bit is used to compute the set size. */
	if st.wMin < 0 {
		st.wMin = index
		st.wordMin = word
	}
	/*  The last non-zero word contains the last actual bit of the set.
	    The location of this bit is used to compute the set length. */
	st.wMax = index
	st.wordMax = word
	/*  Count the actual bits, so as to get the cardinality of the set. */
	st.cardinality = st.cardinality + int32(bits.OnesCount64(word))
}

//-----------------------------------------------------------------------------
/**
 *  The XOR of level3 blocks is computed.
 *
 * <pre>
 * xor| 0 1
 *   0| 0 1
 *   1| 1 0 <pre>
 */
type xorStrategyType struct{}

func (st xorStrategyType) properties() int32 {
	return cFalseOpFalseEqFalse + cValueOpFalseEqValue
}
func (st xorStrategyType) start(b *BitSet) bool {
	return true
}

func (st xorStrategyType) word(base, u3 int32, a3, b3 b1DimType, mask wordType) bool {
	a3[u3] = a3[u3] ^ (b3[u3] & mask)
	return a3[u3] == 0
}

func (st xorStrategyType) block(base, u3, v3 int32, a3, b3 b1DimType) (isZero bool) {
	for w3 := u3; w3 != v3; w3 = w3 + 1 {
		a3[w3] = a3[w3] ^ b3[w3]
		isZero = isZero && a3[w3] == 0
	}
	isZero = false
	return
}
func (st xorStrategyType) finish(cache *cacheType, a2Count, a3Count int32) {}
