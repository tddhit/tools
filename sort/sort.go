package sort

func InsertSort(data []int) {
	if len(data) < 2 {
		return
	}
	for i := 1; i < len(data); i++ {
		for j := i - 1; j >= 0; j-- {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			} else {
				break
			}
		}
	}
}

func QuickSort(data []int) {
	if len(data) < 2 {
		return
	}
	quickSort(data, 0, len(data)-1)
}

func quickSort(data []int, low, high int) {
	if low >= high {
		return
	}
	first, last, key := low, high, data[low]
	for first < last {
		for first < last && data[last] >= key {
			last--
		}
		data[first] = data[last]
		for first < last && data[first] <= key {
			first++
		}
		data[last] = data[first]
	}
	data[first] = key
	quickSort(data, low, first-1)
	quickSort(data, first+1, high)
}

func SelectSort(data []int) {
	for i := 0; i < len(data); i++ {
		min := i
		for j := i + 1; j < len(data); j++ {
			if data[j] < data[min] {
				min = j
			}
		}
		data[i], data[min] = data[min], data[i]
	}
}

func BubbleSort(data []int) {
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func ShellSort(data []int) {
	for gap := len(data) / 2; gap > 0; gap /= 2 {
		for i := gap; i < len(data); i++ {
			j := i
			for j-gap >= 0 && data[j] < data[j-gap] {
				data[j], data[j-gap] = data[j-gap], data[j]
				j -= gap
			}
		}
	}
}

func HeapSort(data []int) {
	buildHeap(data)
	for i := len(data) - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0]
		adjustHeap(data, 0, i-1)
	}
}

func buildHeap(data []int) {
	for i := len(data)/2 - 1; i >= 0; i-- {
		adjustHeap(data, i, len(data)-1)
	}
}

func adjustHeap(data []int, s, n int) {
	max := s
	if 2*s <= n {
		if data[2*s] > data[max] {
			max = 2 * s
		}
	}
	if 2*s+1 <= n {
		if data[2*s+1] > data[max] {
			max = 2*s + 1
		}
	}
	if max != s {
		data[s], data[max] = data[max], data[s]
		adjustHeap(data, max, n)
	}
}

func MergeSort(data []int) {
	if len(data) < 2 {
		return
	}
	res := mergeSort(data)
	copy(data, res)
}

func mergeSort(data []int) []int {
	if len(data) < 2 {
		return data
	}
	m := len(data) / 2
	left := mergeSort(data[:m])
	right := mergeSort(data[m:])
	res := mergeTwo(left, right)
	return res
}

func mergeTwo(left, right []int) (res []int) {
	l, r := 0, 0
	for l < len(left) && r < len(right) {
		if left[l] < right[r] {
			res = append(res, left[l])
			l++
		} else {
			res = append(res, right[r])
			r++
		}
	}
	res = append(res, left[l:]...)
	res = append(res, right[r:]...)
	return
}
