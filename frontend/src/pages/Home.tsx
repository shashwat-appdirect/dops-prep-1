import Hero from '../components/Hero'
import SessionsSpeakers from '../components/SessionsSpeakers'
import RegistrationForm from '../components/RegistrationForm'
import Location from '../components/Location'
import Footer from '../components/Footer'

const Home = () => {
  return (
    <div className="min-h-screen">
      <Hero />
      <SessionsSpeakers />
      <RegistrationForm />
      <Location />
      <Footer />
    </div>
  )
}

export default Home

