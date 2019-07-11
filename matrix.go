package main

type Matrix4x4 [16]float64

func (mat *Matrix4x4) Vec4xMat4(vec [4]float64) [4]float64 {
	return [4]float64{
		mat[0]*vec[0] + mat[1]*vec[0] + mat[2]*vec[0] + mat[3]*vec[0],
		mat[4]*vec[1] + mat[5]*vec[1] + mat[6]*vec[1] + mat[7]*vec[1],
		mat[8]*vec[2] + mat[9]*vec[2] + mat[10]*vec[2] + mat[11]*vec[2],
		mat[12]*vec[3] + mat[13]*vec[3] + mat[14]*vec[3] + mat[15]*vec[3],
	}
}
