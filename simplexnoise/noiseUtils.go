package main

func rescaleAndDraw(rawNoise []float32, minNoise, maxNoise float32, colorGradient []color, pixels []byte) {

	scale := 255.0 / (maxNoise - minNoise)
	offset := minNoise * scale

	for i := range rawNoise {

		rawNoise[i] = rawNoise[i]*scale - offset

		color := colorGradient[clamp(0, 255, int(rawNoise[i]))]

		pixels[i*channelDepth] = color.red
		pixels[i*channelDepth+1] = color.green
		pixels[i*channelDepth+2] = color.blue
	}
}

func fractionalBrownianMotion(x, y, frequency, lacunarity, gain float32, octaves int) float32 {

	var sum float32

	amplitude := float32(1.0)

	for i := 0; i < octaves; i++ {

		sum += snoise2(x*frequency, y*frequency) * amplitude

		frequency *= lacunarity
		amplitude *= gain
	}

	return sum
}

func makeNoise(pixels []byte, frequency, gain, lacunarity float32, octaves int) {

	noise := make([]float32, windowHeight*windowWidth)

	i := 0

	var minNoise float32
	var maxNoise float32

	for y := 0; y < windowWidth; y++ {
		for x := 0; x < windowHeight; x++ {
			noise[i] = fractionalBrownianMotion(float32(x), float32(y), frequency, lacunarity, gain, octaves)
			if i == 0 {
				minNoise = noise[i]
				maxNoise = noise[i]
			} else if noise[i] > maxNoise {
				maxNoise = noise[i]
			} else if noise[i] < minNoise {
				minNoise = noise[i]
			}
			i++
		}
	}
	gradient := getColorGradient(color{255, 0, 0}, color{0, 255, 0})
	rescaleAndDraw(noise, minNoise, maxNoise, gradient, pixels)
}

func lerp(b1 byte, b2 byte, percent float32) byte {
	return byte(float32(b1) + percent*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, percent float32) color {
	return color{lerp(c1.red, c2.red, percent), lerp(c1.blue, c2.blue, percent), lerp(c1.green, c2.green, percent)}
}

func getColorGradient(c1, c2 color) []color {

	colorGradient := make([]color, 256)

	for i := range colorGradient {
		percent := float32(i) / float32(255)
		colorGradient[i] = colorLerp(c1, c2, percent)
	}

	return colorGradient
}

func clamp(min, max, value int) int {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	return value
}