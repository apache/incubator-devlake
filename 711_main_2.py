def quick_sort(arr):
    """快速排序算法"""
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]  # 选择中间元素作为基准
    left = [x for x in arr if x < pivot]  # 小于基准的元素
    middle = [x for x in arr if x == pivot]  # 等于基准的元素
    right = [x for x in arr if x > pivot]  # 大于基准的元素
    return quick_sort(left) + middle + quick_sort(right)  # 递归排序并合并

# 示例用法
arr = [3, 6, 8, 10, 1, 2, 1]
print("未排序的数组:", arr)
sorted_arr = quick_sort(arr)
print("已排序的数组:", sorted_arr)
# zzehegg
@zehgsjhshvshwkk
#hhhhhhhhhh

@hshdjsjfhs
#hhsdjshdjshdhfsh
