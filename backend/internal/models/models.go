package models

import "time"

type Registration struct {
	ID          string    `json:"id" firestore:"-"`
	Name        string    `json:"name" firestore:"name"`
	Email       string    `json:"email" firestore:"email"`
	Designation string    `json:"designation" firestore:"designation"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
}

type Speaker struct {
	ID          string `json:"id" firestore:"-"`
	Name        string `json:"name" firestore:"name"`
	Bio         string `json:"bio" firestore:"bio"`
	ImageURL    string `json:"imageUrl,omitempty" firestore:"imageUrl,omitempty"`
	LinkedInURL string `json:"linkedinUrl,omitempty" firestore:"linkedinUrl,omitempty"`
	TwitterURL  string `json:"twitterUrl,omitempty" firestore:"twitterUrl,omitempty"`
}

type Session struct {
	ID          string   `json:"id" firestore:"-"`
	Title       string   `json:"title" firestore:"title"`
	Description string   `json:"description" firestore:"description"`
	Time        string   `json:"time" firestore:"time"`
	Duration    string   `json:"duration" firestore:"duration"`
	SpeakerIDs  []string `json:"speakerIds" firestore:"speakerIds"`
}

type DesignationBreakdown struct {
	Designation string `json:"designation"`
	Count       int    `json:"count"`
}

