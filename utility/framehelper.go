package utility

// IsFrameInPast ...
func IsFrameInPast(frame byte, currentFrame byte) bool {
	return currentFrame > frame || (currentFrame >= 0 && currentFrame < 30 && (255-frame) < 30)
}

// IsFrameNowOrInPast ...
func IsFrameNowOrInPast(frame byte, currentFrame byte) bool {
	return currentFrame >= frame || (currentFrame >= 0 && currentFrame < 30 && (255-frame) < 30)
}
