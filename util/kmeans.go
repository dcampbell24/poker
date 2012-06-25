package util

func sort(a []float64) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j] < a[j-1]; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}

func resizef(x []float64, y int) []float64 {
	if cap(x) < y {
		z := make([]float64, y)
		copy(z, x)
		return z
	}
	return x[:y]
}

func resize(x []int, l int) []int {
	if cap(x) < l {
		z := make([]int, l)
		copy(z, x)
		return z
	}
	return x[:l]
}

/*
Ckmeans_1d_dp.h --- Head file for Ckmeans.1d.dp
                    Declare a class "data" and 
					wrap function "kmeans_1d_dp()"

  Haizhou Wang
  Computer Science Department
  New Mexico State University
  hwang@cs.nmsu.edu

Created: Oct 10, 2010
*/

// data that return by kmeans.1d.dp()
type data struct {
	Cluster  []int     // record which cluster each point belongs to
	Centers  []float64 // record the center of each cluster
	WithinSS []float64 // within sum of distance square of each cluster
	Size     []int     // size of each cluster
}

// one-dimensional cluster algorithm implemented in Go
// x is input one-dimensional vector and K stands for the cluster level
// data kmeans_1d_dp( vector<double> x, int K);

// END HEADER
/*
Ckmeans_1d_dp.cpp -- Performs 1-D k-means by a dynamic programming
                     approach that is guaranteed to be optimal.

  Joe Song
  Computer Science Department
  New Mexico State University
  joemsong@cs.nmsu.edu

  Haizhou Wang
  Computer Science Department
  New Mexico State University
  hwang@cs.nmsu.edu

Created: May 19, 2007
Updated: September 3, 2009
Updated: September 6, 2009.  Handle special cases when K is too big or array
	 contains identical elements.  Added number of clusters selection by the MCLUST package.
Updated: Oct 6, 2009 by Haizhou Wang Convert this program from R code into C++ code.
Updated: Oct 20, 2009 by Haizhou Wang, add interface to this function, so that it could be called directly in R
*/

// all vectors in this program is considered starting at position 1, position 0 is not used.
//
// Input:
//	x -- a vector of numbers, not necessarily sorted
//	K -- the number of clusters expected
//
// Pre-Conditions:
//	K <= |set(x)|
//	|set(x)| > 1
func kmeans_1d_dp(x []float64, K int) data {
	result := data{}
	N := len(x) - 1 // N: is the size of input vector
	y := make([]int, len(x))
	temp_s := make([]float64, len(x))
	copy(temp_s, x)

	// create a mapping from the unsorted to the sorted order of the input.
	sort(temp_s[1:])
	for i := 1; i < len(x); i++ {
		for j := 1; j < len(x); j++ {
			if x[i] == temp_s[j] {
				y[i] = j
				break
			}
		}
	}
	sort(x[1:])
	D := make([][]float64, K+1)
	for i := range D {
		D[i] = make([]float64, N+1)
	}
	B := make([][]float64, K+1)
	for i := range B {
		B[i] = make([]float64, N+1)
	}
	for i := 1; i <= K; i++ {
		D[i][1] = 0
		B[i][1] = 1
	}
	var mean_x1, mean_xj, d float64
	for k := 1; k <= K; k++ {
		mean_x1 = x[1]
		for i := 2; i <= N; i++ {
			if k == 1 {
				D[1][i] = D[1][i-1] + float64(i-1)/float64(i)*(x[i]-mean_x1)*(x[i]-mean_x1)
				mean_x1 = (float64(i-1)*mean_x1 + x[i])/float64(i)
				B[1][i] = 1
			} else {
				D[k][i] = -1
				d = 0
				mean_xj = 0
				for j := i; j >= 1; j-- {
					d += float64(i - j) / float64(i-j+1) * (x[j] - mean_xj) * (x[j] - mean_xj)
					mean_xj = (x[j] + float64(i-j)*mean_xj) / float64(i-j+1)
					//initialization of D[k,i]
					if D[k][i] == -1 {
						if j == 1 {
							D[k][i] = d
							B[k][i] = float64(j)
						} else {
							D[k][i] = d + D[k-1][j-1]
							B[k][i] = float64(j)
						}
					} else {
						if j == 1 {
							if d <= D[k][i] {
								D[k][i] = d
								B[k][i] = float64(j)
							}
						} else {
							if d+D[k-1][j-1] < D[k][i] {
								D[k][i] = d + D[k-1][j-1]
								B[k][i] = float64(j)
							}
						}
					}
				}
			}
		}
	}
	// Backtrack to find the clusters of the data points
	cluster_right := N
	var cluster_left float64
	result.Cluster = resize(result.Cluster, N+1)
	result.Centers = resizef(result.Centers, K+1)
	result.WithinSS = resizef(result.WithinSS, K+1)
	result.Size = resize(result.Size, K+1)

	/*Forming final result*/
	for k := K; k >= 1; k-- {
		cluster_left = B[k][cluster_right]
		for i := int(cluster_left); i <= cluster_right; i++ {
			result.Cluster[i] = k
		}
		var sum float64
		for a := int(cluster_left); a <= cluster_right; a++ {
			sum += x[a]
		}
		result.Centers[k] = sum / (float64(cluster_right) - cluster_left + 1)
		for a := int(cluster_left); a <= cluster_right; a++ {
			result.WithinSS[k] += (x[a] - result.Centers[k]) * (x[a] - result.Centers[k])
		}
		result.Size[k] = cluster_right - int(cluster_left) + 1
		if k > 1 {
			cluster_right = int(cluster_left) - 1
		}
	}
	// restore the original order
	tt := result.Cluster
	for i := 1; i < len(x); i++ {
		result.Cluster[i] = tt[y[i]]
	}
	return result
}
