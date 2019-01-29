package sparse

import "fmt"

//=============================================================================
//  Statistics enumeration
//=============================================================================

/**
 *  These enumeration values are used as labels for the values in the String
 *  created by the <i>statistics</i>() method. The values of the corresponding
 *  statistics are <code>int</code>s, except for the loadFactor and
 *  Average_chain_length values, which are <code>float</code>s.
 *  <p>
 *  An array of <code>String</code>s may be obtained containing a
 *  representation of each of these values. An element of such an array, say,
 *  <code>values</code>, may be accessed, for example, by:
 *  <pre>
 *      values[SparseBitSet.statistics.Buckets_available.ordinal()]</pre>
 *
 * @see         #statistics(String[])
 */
type StatisticsType int

const (
	/**
	 *  The size of the bit set, as give by the <i>size</i>() method.
	 */
	Size StatisticsType = iota // 0
	/**
	 *  The length of the bit set, as give by the <i>length</i>() method.
	 */
	Length // 1
	/**
	 *  The cardinality of the bit set, as give by the <i>cardinality</i>() method.
	 */
	Cardinality // 2
	/**
	 *  The total number of non-zero 64-bits "words" being used to hold the
	 *  representation of the bit set.
	 */
	Total_words // 3
	/**
	 *  The length of the bit set array.
	 */
	Set_array_length // 4
	/**
	 *  The maximum permitted length of the bit set array.
	 */
	Set_array_max_length // 5
	/**
	 *  The number of level2 areas.
	 */
	Level2_areas // 6
	/**
	 *  The length of the level2 areas.
	 */
	Level2_area_length // 7
	/**
	 *  The total number of level3 blocks in use.
	 */
	Level3_blocks // 8
	/**
	 *  The length of the level3 blocks.
	 */
	Level3_block_length // 9
	/**
	 *  Is the value that determines how the <i>toString</i>() conversion is
	 *  performed.
	 * @see         #toStringCompaction(int)
	 */
	Compaction_count_value // 10
	//
	Statistics_Values_Length
)

func (st StatisticsType) String() string {
	switch st {
	case Size:
		return "Size"
	case Length:
		return "Length"
	case Cardinality:
		return "Cardinality"
	case Total_words:
		return "Total-words"
	case Set_array_length:
		return "Set-array-length"
	case Set_array_max_length:
		return "Set-array-max-length"
	case Level2_areas:
		return "Level2-areas"
	case Level2_area_length:
		return "Level2-area-length"
	case Level3_blocks:
		return "Level3-blocks"
	case Level3_block_length:
		return "Level3-block-length"
	case Compaction_count_value:
		return "Compaction-count-value"
	default:
		panic(fmt.Sprintf("Unknown statistics value %d", st))
	}
}
