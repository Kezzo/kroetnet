package collision

import "testing"

func TestCircleCollisionSimple(t *testing.T) {
	cc := CircleCollider{0, 0, 5}
	cc2 := CircleCollider{2, 3, 5}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestCircleCollisionOnSamePosition(t *testing.T) {
	cc := CircleCollider{0, 0, 5}
	cc2 := CircleCollider{0, 0, 5}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestCircleCollisionTangentTo(t *testing.T) {
	cc := CircleCollider{0, 0, 5}
	cc2 := CircleCollider{10, 0, 5}

	if !cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestCircleCollisionNotColliding(t *testing.T) {
	cc := CircleCollider{0, 0, 5}
	cc2 := CircleCollider{12, 15, 5}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestCircleCollisionUnsupportedCollider(t *testing.T) {
	cc := CircleCollider{0, 0, 5}
	cc2 := ConeCollider{6, 6, 5, 90, 45}

	if cc.IsColliding(&cc2) {
		t.Errorf("Collision was not correctly detected")
	}
}

func TestCircleColliderPositionUpdates(t *testing.T) {
	cc := CircleCollider{0, 0, 5}

	cc.Update(12, 15, 127)

	if cc.Xpos != 12 || cc.Ypos != 15 || cc.Radius != 5 {
		t.Errorf("CircleCollider position was not correctly updated")
	}
}
