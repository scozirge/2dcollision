package collision

import (
	"github.com/scozirge/2dcollision/collider"
	"math"
)

// 判斷圓形與對齊兩軸的矩形之間是否有碰撞
func AxisAligned_CircleToRectangleColliding(c collider.Circle, r collider.Rectangle1) bool {
	// 從圓心找到最接近的矩形的x點
	closestX := clamp(c.X, r.X, r.X+r.Width)
	// 從圓心找到最接近的矩形的y點
	closestY := clamp(c.Y, r.Y, r.Y+r.Height)

	// 計算圓心與這最接近的點之間的距離的平方
	distanceX := c.X - closestX
	distanceY := c.Y - closestY
	distanceSquared := distanceX*distanceX + distanceY*distanceY

	// 以下使用距離平方與半徑平方來比較，避免開方在計算中的開銷
	// 如果距離的平方小於圓的半徑的平方，則發生交集
	return distanceSquared <= c.Radius*c.Radius
}

// 判斷圓形與旋轉的矩形間是否有碰撞(使用分離軸定理 Separating Axis Theorem, SAT)
func CircleToRectangleColliding(c collider.Circle, r collider.Rectangle2) bool {
	axes := []collider.Vector2{
		{r.P1.X - r.P2.X, r.P1.Y - r.P2.Y},
		{r.P1.X - r.P4.X, r.P1.Y - r.P4.Y},
		{c.X - r.P1.X, c.Y - r.P1.Y},
		{c.X - r.P2.X, c.Y - r.P2.Y},
		{c.X - r.P3.X, c.Y - r.P3.Y},
		{c.X - r.P4.X, c.Y - r.P4.Y},
	}

	for _, axis := range axes {
		if !overlap(projectRectangle(r, axis), projectCircle(c, axis)) {
			return false
		}
	}

	return true
}

// 將矩形投影到軸上
func projectRectangle(rect collider.Rectangle2, axis collider.Vector2) (float64, float64) {
	dots := []float64{
		dot(axis, rect.P1),
		dot(axis, rect.P2),
		dot(axis, rect.P3),
		dot(axis, rect.P4),
	}

	min := dots[0]
	max := dots[0]
	for _, d := range dots[1:] {
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	return min, max
}

// 將圓形投影到軸上
func projectCircle(c collider.Circle, axis collider.Vector2) (float64, float64) {
	projectionCenter := dot(axis, collider.Vector2{X: c.X, Y: c.Y})
	return projectionCenter - c.Radius, projectionCenter + c.Radius
}

// 檢查兩個投影是否有重疊
func overlap(minMax1, minMax2 (float64, float64)) bool {
	return minMax1[1] >= minMax2[0] && minMax2[1] >= minMax1[0]
}

// 計算兩個向量的點積
func dot(v1, v2 collider.Vector2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}
// 限制值在最小和最大之間
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
