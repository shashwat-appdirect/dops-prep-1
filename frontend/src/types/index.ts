export interface Registration {
  id?: string
  name: string
  email: string
  designation: string
  createdAt?: string
}

export interface Speaker {
  id?: string
  name: string
  bio: string
  imageUrl?: string
  linkedinUrl?: string
  twitterUrl?: string
}

export interface Session {
  id?: string
  title: string
  description: string
  time: string
  duration: string
  speakerIds: string[]
}

export interface SessionWithSpeakers extends Session {
  speakers: Speaker[]
}

export interface DesignationBreakdown {
  designation: string
  count: number
}

