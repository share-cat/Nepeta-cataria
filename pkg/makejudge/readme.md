
## 類型映射表

| 類型描述 | Python | C++ | JavaScript |
|---------|--------|-----|------------|
| `int` | `int` | `int` | `number` |
| `str` | `str` | `std::string` | `string` |
| `float` | `float` | `float` | `number` |
| `double` | `float` | `double` | `number` |
| `list[int]` | `List[int]` | `std::vector<int>` | `number[]` |
| `list[str]` | `List[str]` | `std::vector<std::string>` | `string[]` |
| `list[list[int]]` | `List[List[int]]` | `std::vector<std::vector<int>>` | `number[][]` |
| `set[int]` | `Set[int]` | `std::set<int>` | `Set<number>` |
| `dict[int]` | `Dict[str, int]` | `std::map<std::string, int>` | `Record<string, number>` |

## testdata.yaml 配置

### 範例 1: 雨水收集問題 (二維矩陣)

```yaml
config:
  input_format:
    type: "list[list[int]]"
  output_format:
    type: "int"

sample:
  input: |
    [[1,4,3],[3,2,1],[2,3,3]]
  output: |
    4
```

**生成的 Python 模板：**
```python
from typing import List

def main(input_data: List[List[int]]) -> int:
    # Replace this with your solution.
    raise NotImplementedError("Please implement the solution.")

```

**生成的 C++ 模板：**
```cpp
#include <iostream>
#include <sstream>
#include <stdexcept>
#include <vector>
#include <string>

int solve(const std::vector<std::vector<int>>& input_data) {
    // Replace this with your solution.
    throw std::runtime_error("Please implement the solution.");
}
```

### 範例 2: 字符串處理問題

```yaml
config:
  input_format:
    type: "list[str]"
  output_format:
    type: "str"

sample:
  input: |
    ["hello", "world", "test"]
  output: |
    "hello world test"
```

**生成的 Python 函數：**
```python
def main(input_data: List[str]) -> str:
```

**生成的 C++ 函數：**
```cpp
std::string solve(const std::vector<std::string>& input_data) {
```

### 範例 3: 複雜數據結構

```yaml
config:
  input_format:
    type: "list[list[list[int]]]"  # 三維陣列
  output_format:
    type: "list[int]"              # 一維陣列

sample:
  input: |
    [[[1,2],[3,4]], [[5,6],[7,8]]]
  output: |
    [10, 26]
```

**生成的類型：**
- **Python**: `List[List[List[int]]]` → `List[int]`
- **C++**: `std::vector<std::vector<std::vector<int>>>` → `std::vector<int>`
- **JavaScript**: `number[][][]` → `number[]`

## 使用流程

### 1. 創建問題
```bash
mkdir -p judges/matrix_problem
cd judges/matrix_problem
```

### 2. 編寫 config.yaml
```yaml
name: "Matrix Rain Water"
description: "計算矩陣中的雨水收集量"
lang:
  - python
  - cpp
  - javascript
level: medium
type: judge
mem: 256m
time: 2s
```

### 3. 編寫 testdata.yaml
```yaml
config:
  input_format:
    type: "list[list[int]]"
    description: "2D integer matrix"
  output_format:
    type: "int"
    description: "Total water trapped"

sample:
  input: |
    [[1,4,3],[3,2,1],[2,3,3]]
  output: |
    4

test_cases:
  - name: "test_01"
    input: |
      [[1,4,3],[3,2,1],[2,3,3]]
    output: |
      4
    points: 50
```

### 4. 生成代碼模板
```bash
# 生成 Python 模板
makejudge generate -problem_dir judges/matrix_problem -lang python

# 生成 C++ 模板
makejudge generate -problem_dir judges/matrix_problem -lang cpp
```
