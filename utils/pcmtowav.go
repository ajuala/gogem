// Package audio provides utilities for audio manipulation, including PCM to WAV conversion.
package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	// "os" // Uncomment for example usage
)

// WavHeader represents the structure of a standard WAV file header (RIFF, fmt, and data chunks).
// This struct is designed to be directly written using binary.Write with Little Endian byte order.
type WavHeader struct {
	// RIFF Chunk
	RiffID        [4]byte // Contains "RIFF"
	ChunkSize     uint32  // Size of the rest of the file after this field (ChunkSize + WAVE + fmt chunk + data chunk)
	WaveID        [4]byte // Contains "WAVE"

	// fmt Sub-chunk
	FmtID         [4]byte // Contains "fmt " (note the trailing space)
	Subchunk1Size uint32  // Size of the fmt sub-chunk (16 for PCM)
	AudioFormat   uint16  // Audio format (1 for PCM)
	NumChannels   uint16  // Number of channels (1 for mono, 2 for stereo)
	SampleRate    uint32  // Sample rate (samples per second, e.g., 44100)
	ByteRate      uint32  // Byte rate (SampleRate * NumChannels * BitsPerSample/8)
	BlockAlign    uint16  // Block align (NumChannels * BitsPerSample/8)
	BitsPerSample uint16  // Bits per sample (e.g., 8, 16, 24, 32)

	// Data Sub-chunk
	DataID        [4]byte // Contains "data"
	Subchunk2Size uint32  // Size of the data sub-chunk (number of bytes in the PCM data)
}

// ConvertPCMToWav converts raw PCM audio data into a WAV formatted byte slice.
//
// Parameters:
//   pcmData       []byte: The raw PCM audio data bytes. This data should be
//                         interleaved if multiple channels are present (e.g.,
//                         left sample, right sample, left sample, right sample for stereo).
//                         The byte order within samples should match the system's
//                         endianness or be consistent (typically Little Endian for WAV).
//   numChannels   int: The number of audio channels (e.g., 1 for mono, 2 for stereo).
//   sampleRate    int: The sample rate in Hz (e.g., 44100).
//   bitsPerSample int: The number of bits per audio sample (e.g., 8, 16, 24, 32).
//
// Returns:
//   []byte: A byte slice containing the complete WAV file data.
//   error: An error if conversion fails (e.g., invalid parameters or writing issues).
func ConvertPCMToWav(pcmData []byte, numChannels, sampleRate, bitsPerSample int) ([]byte, error) {
	// Validate input parameters
	if numChannels <= 0 {
		return nil, fmt.Errorf("number of channels must be positive, got %d", numChannels)
	}
	if sampleRate <= 0 {
		return nil, fmt.Errorf("sample rate must be positive, got %d", sampleRate)
	}
	if bitsPerSample%8 != 0 || bitsPerSample <= 0 {
		return nil, fmt.Errorf("bits per sample must be a multiple of 8 and positive, got %d", bitsPerSample)
	}

	bytesPerSample := bitsPerSample / 8
	// Calculate blockAlign: The number of bytes for one frame (all channels at one sample point).
	blockAlign := uint16(numChannels * bytesPerSample)
	// Calculate byteRate: The average number of bytes per second.
	byteRate := uint32(sampleRate * int(blockAlign))
	// subchunk2Size: The size of the actual PCM data chunk.
	subchunk2Size := uint32(len(pcmData))

	// chunkSize: Total file size - 8 bytes (for "RIFF" and ChunkSize fields themselves).
	// It's 4 (for "WAVE") + (8 + Subchunk1Size) + (8 + Subchunk2Size).
	// For PCM, Subchunk1Size is 16. So, 4 + 8 + 16 + 8 + subchunk2Size = 36 + subchunk2Size.
	chunkSize := uint32(36 + subchunk2Size)

	// Populate the WAV header struct
	header := WavHeader{
		RiffID:        [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     chunkSize,
		WaveID:        [4]byte{'W', 'A', 'V', 'E'},
		FmtID:         [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16, // For PCM
		AudioFormat:   1,  // For PCM (Pulse Code Modulation)
		NumChannels:   uint16(numChannels),
		SampleRate:    uint32(sampleRate),
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: uint16(bitsPerSample),
		DataID:        [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: subchunk2Size,
	}

	// Create a buffer to write the WAV data to
	buf := new(bytes.Buffer)

	// Write the WAV header to the buffer.
	// WAV files use Little Endian byte order.
	if err := binary.Write(buf, binary.LittleEndian, header); err != nil {
		return nil, fmt.Errorf("failed to write WAV header: %w", err)
	}

	// Write the raw PCM data after the header.
	if err := binary.Write(buf, binary.LittleEndian, pcmData); err != nil {
		return nil, fmt.Errorf("failed to write PCM data: %w", err)
	}

	return buf.Bytes(), nil
}

/*
// Example Usage (uncomment to run)
func main() {
	// --- Example 1: Generate a simple 440 Hz sine wave (mono, 16-bit, 44100 Hz) ---
	sampleRate := 44100
	frequency := 440.0 // Hz
	duration := 3     // seconds
	amplitude := 0.5  // Max amplitude (for float32, converted to int16)
	bitsPerSample := 16
	numChannels := 1

	// Calculate number of samples
	numSamples := sampleRate * duration
	pcmData := make([]byte, numSamples*numChannels*(bitsPerSample/8))

	// Generate sine wave samples
	for i := 0; i < numSamples; i++ {
		// Calculate sample value (float between -1.0 and 1.0)
		t := float64(i) / float64(sampleRate)
		sampleValue := amplitude * math.Sin(2*math.Pi*frequency*t)

		// Convert float to int16 (PCM 16-bit ranges from -32768 to 32767)
		// and convert to little-endian bytes.
		int16Sample := int16(sampleValue * 32767.0)
		binary.LittleEndian.PutUint16(pcmData[i*2:], uint16(int16Sample))
	}

	// Convert PCM data to WAV format
	wavBytes, err := ConvertPCMToWav(pcmData, numChannels, sampleRate, bitsPerSample)
	if err != nil {
		fmt.Println("Error converting to WAV:", err)
		return
	}

	// Save to a WAV file
	err = os.WriteFile("sine_wave.wav", wavBytes, 0644)
	if err != nil {
		fmt.Println("Error writing WAV file:", err)
		return
	}
	fmt.Println("Generated sine_wave.wav successfully!")


	// --- Example 2: Convert existing raw PCM file (if you have one) ---
	// Assume 'input.pcm' is a raw 16-bit, mono, 44100 Hz PCM file
	// inputPCMFilePath := "input.pcm"
	// outputWAVFilePath := "output.wav"
	// pcmDataFromFile, err := os.ReadFile(inputPCMFilePath)
	// if err != nil {
	// 	fmt.Printf("Error reading PCM file %s: %v\n", inputPCMFilePath, err)
	// 	return
	// }
	//
	// wavBytesFromFile, err := ConvertPCMToWav(pcmDataFromFile, 1, 44100, 16)
	// if err != nil {
	// 	fmt.Println("Error converting raw PCM file to WAV:", err)
	// 	return
	// }
	//
	// err = os.WriteFile(outputWAVFilePath, wavBytesFromFile, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing output WAV file:", err)
	// 	return
	// }
	// fmt.Printf("Converted %s to %s successfully!\n", inputPCMFilePath, outputWAVFilePath)
}

// import "math" // Uncomment for example usage
*/
