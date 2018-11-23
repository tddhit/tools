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
