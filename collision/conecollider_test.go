package collision

import "testing"

func TestConeCollisionSimple45(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 45}
	cc2 := CircleCollider{2, 3, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionSimple45(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 45}
	cc2 := CircleCollider{2, -3, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionSimple90(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 90}
	cc2 := CircleCollider{5, 1, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionSimple90(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 90}
	cc2 := CircleCollider{5, -2, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionSimple135(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 135}
	cc2 := CircleCollider{5, -2, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionSimple135(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 135}
	cc2 := CircleCollider{-1, -4, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionRotated90(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 90, 45}
	cc2 := CircleCollider{3, -1, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionRotated90(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 90, 45}
	cc2 := CircleCollider{-3, -1, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionRotated180(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 180, 45}
	cc2 := CircleCollider{1, -3, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionRotated180(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 180, 45}
	cc2 := CircleCollider{-3, 1, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionRotated270(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 270, 45}
	cc2 := CircleCollider{-3, 1, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionRotated270(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 270, 45}
	cc2 := CircleCollider{1, -3, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionRotated225(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 225, 45}
	cc2 := CircleCollider{-3, -3, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeNoCollisionRotated225(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 225, 45}
	cc2 := CircleCollider{1, -3, 1}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionNoDegressDifference(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 45}
	cc2 := CircleCollider{0, 4, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionRadiusTouchedCone(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 45}
	cc2 := CircleCollider{4, 4, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestConeCollisionSamePosition(t *testing.T) {
	cc := ConeCollider{0, 0, 5, 0, 45}
	cc2 := CircleCollider{0, 0, 1}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}
