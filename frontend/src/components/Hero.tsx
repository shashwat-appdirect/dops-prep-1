const Hero = () => {
  return (
    <section className="gradient-bg text-white py-20 md:py-32 px-4">
      <div className="max-w-6xl mx-auto text-center animate-fade-in">
        <h1 className="text-4xl md:text-6xl font-bold mb-6 animate-slide-down">
          AppDirect India AI Workshop
        </h1>
        <p className="text-xl md:text-2xl mb-8 text-gray-100 max-w-3xl mx-auto animate-slide-up">
          Join us for an exciting day of AI innovation, expert sessions, and networking opportunities.
          Discover the future of artificial intelligence and connect with industry leaders.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center animate-slide-up">
          <a
            href="#register"
            className="bg-white text-blue-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100 transition-all duration-300 transform hover:scale-105 shadow-lg"
          >
            Register Now
          </a>
          <a
            href="#sessions"
            className="bg-transparent border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white hover:text-blue-600 transition-all duration-300 transform hover:scale-105"
          >
            View Sessions
          </a>
        </div>
      </div>
    </section>
  )
}

export default Hero

