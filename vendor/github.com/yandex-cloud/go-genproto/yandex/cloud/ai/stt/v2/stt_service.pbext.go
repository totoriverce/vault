// Code generated by protoc-gen-goext. DO NOT EDIT.

package stt

import (
	duration "github.com/golang/protobuf/ptypes/duration"
)

func (m *LongRunningRecognitionRequest) SetConfig(v *RecognitionConfig) {
	m.Config = v
}

func (m *LongRunningRecognitionRequest) SetAudio(v *RecognitionAudio) {
	m.Audio = v
}

func (m *LongRunningRecognitionResponse) SetChunks(v []*SpeechRecognitionResult) {
	m.Chunks = v
}

type StreamingRecognitionRequest_StreamingRequest = isStreamingRecognitionRequest_StreamingRequest

func (m *StreamingRecognitionRequest) SetStreamingRequest(v StreamingRecognitionRequest_StreamingRequest) {
	m.StreamingRequest = v
}

func (m *StreamingRecognitionRequest) SetConfig(v *RecognitionConfig) {
	m.StreamingRequest = &StreamingRecognitionRequest_Config{
		Config: v,
	}
}

func (m *StreamingRecognitionRequest) SetAudioContent(v []byte) {
	m.StreamingRequest = &StreamingRecognitionRequest_AudioContent{
		AudioContent: v,
	}
}

func (m *StreamingRecognitionResponse) SetChunks(v []*SpeechRecognitionChunk) {
	m.Chunks = v
}

type RecognitionAudio_AudioSource = isRecognitionAudio_AudioSource

func (m *RecognitionAudio) SetAudioSource(v RecognitionAudio_AudioSource) {
	m.AudioSource = v
}

func (m *RecognitionAudio) SetContent(v []byte) {
	m.AudioSource = &RecognitionAudio_Content{
		Content: v,
	}
}

func (m *RecognitionAudio) SetUri(v string) {
	m.AudioSource = &RecognitionAudio_Uri{
		Uri: v,
	}
}

func (m *RecognitionConfig) SetSpecification(v *RecognitionSpec) {
	m.Specification = v
}

func (m *RecognitionConfig) SetFolderId(v string) {
	m.FolderId = v
}

func (m *RecognitionSpec) SetAudioEncoding(v RecognitionSpec_AudioEncoding) {
	m.AudioEncoding = v
}

func (m *RecognitionSpec) SetSampleRateHertz(v int64) {
	m.SampleRateHertz = v
}

func (m *RecognitionSpec) SetLanguageCode(v string) {
	m.LanguageCode = v
}

func (m *RecognitionSpec) SetProfanityFilter(v bool) {
	m.ProfanityFilter = v
}

func (m *RecognitionSpec) SetModel(v string) {
	m.Model = v
}

func (m *RecognitionSpec) SetPartialResults(v bool) {
	m.PartialResults = v
}

func (m *RecognitionSpec) SetSingleUtterance(v bool) {
	m.SingleUtterance = v
}

func (m *RecognitionSpec) SetAudioChannelCount(v int64) {
	m.AudioChannelCount = v
}

func (m *RecognitionSpec) SetRawResults(v bool) {
	m.RawResults = v
}

func (m *SpeechRecognitionChunk) SetAlternatives(v []*SpeechRecognitionAlternative) {
	m.Alternatives = v
}

func (m *SpeechRecognitionChunk) SetFinal(v bool) {
	m.Final = v
}

func (m *SpeechRecognitionChunk) SetEndOfUtterance(v bool) {
	m.EndOfUtterance = v
}

func (m *SpeechRecognitionResult) SetAlternatives(v []*SpeechRecognitionAlternative) {
	m.Alternatives = v
}

func (m *SpeechRecognitionResult) SetChannelTag(v int64) {
	m.ChannelTag = v
}

func (m *SpeechRecognitionAlternative) SetText(v string) {
	m.Text = v
}

func (m *SpeechRecognitionAlternative) SetConfidence(v float32) {
	m.Confidence = v
}

func (m *SpeechRecognitionAlternative) SetWords(v []*WordInfo) {
	m.Words = v
}

func (m *WordInfo) SetStartTime(v *duration.Duration) {
	m.StartTime = v
}

func (m *WordInfo) SetEndTime(v *duration.Duration) {
	m.EndTime = v
}

func (m *WordInfo) SetWord(v string) {
	m.Word = v
}

func (m *WordInfo) SetConfidence(v float32) {
	m.Confidence = v
}
