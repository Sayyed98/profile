import { useEffect, useMemo, useState } from "react";
import { fetchPortfolio, formatRange, sendContact } from "./api.js";

export default function App() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");
  const [form, setForm] = useState({ name: "", email: "", subject: "", message: "" });

  useEffect(() => {
    fetchPortfolio()
      .then(setData)
      .finally(() => setLoading(false));
      // .catch(() => setError("Could not reach the Go API. Start MySQL + the backend, then refresh."));
  }, []);

  const years = useMemo(() => {
    if (!data?.experience?.length) return "3+";
    const starts = data.experience.map((e) => new Date(e.start_date).getFullYear());
    return String(new Date().getFullYear() - Math.min(...starts)) + "+";
  }, [data]);

  async function onSubmit(e) {
    e.preventDefault();
    setStatus("Sending…");
    try {
      await sendContact(form);
      setStatus("Message accepted. A worker will process the event asynchronously.");
      setForm({ name: "", email: "", subject: "", message: "" });
    } catch (err) {
      setStatus(err.message);
    }
  }

  return (
    <div className="page">
      <div className="grid-bg" />
      <header className="nav">
        <span className="logo">Mohd Hujaifa/profile</span>
        <nav>
          <a href="#skills">Skills</a>
          <a href="#experience">Experience</a>
          <a href="#projects">Projects</a>
          <a href="#contact">Contact</a>
        </nav>
      </header>

      {error && <div className="banner">{error}</div>}
      {loading && (
        <div className="banner">
            Waking up the backend… first load can take up to a minute (free server tier).
        </div>
      )}

      <section className="hero">
        <p className="kicker">Allahabad · Backend · Distributed systems</p>
        <h1>
          {data?.profile?.full_name || "Mohd Hujaifa"}
          <span>builds fast Golang services.</span>
        </h1>
        <p className="lede">{data?.profile?.summary}</p>
        <div className="hero-meta">
          <a href={`mailto:${data?.profile?.email || "huzaifaanis40@gmail.com"}`}>
            {data?.profile?.email || "huzaifaanis40@gmail.com"}
          </a>
          <span>{data?.profile?.phone}</span>
          <span>{data?.profile?.location}</span>
        </div>
        <div className="stats">
          <Stat value={years} label="Years shipping backends" />
          <Stat value="40%" label="API latency gain via MySQL + Redis" />
          <Stat value="7" label="Microservices on e-commerce" />
        </div>
      </section>

      <section id="skills">
        <h2>Skills</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : (
          <div className="skill-grid">
            {(data?.skills || []).map((g) => (
              <article key={g.category} className="card">
                <h3>{g.category}</h3>
                <ul className="chips">
                  {g.items.map((item) => (
                    <li key={item}>{item}</li>
                  ))}
                </ul>
              </article>
            ))}
          </div>
        )}
      </section>

      <section id="experience">
        <h2>Experience</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : (
          <ol className="timeline">
            {(data?.experience || []).map((job) => (
              <li key={job.id} className="card">
                <div className="job-head">
                  <div>
                    <h3>{job.role}</h3>
                    <p className="muted">
                      {job.company}
                      {job.client ? ` · Client: ${job.client}` : ""} · {job.location}
                    </p>
                  </div>
                  <time>{formatRange(job.start_date, job.end_date)}</time>
                </div>
                <ul className="bullets">
                  {job.bullets.map((b) => (
                    <li key={b}>{b}</li>
                  ))}
                </ul>
              </li>
            ))}
          </ol>
        )}
      </section>

      <section id="projects">
        <h2>Projects</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : (
          <div className="project-grid">
            {(data?.projects || []).map((p) => (
              <article key={p.id} className="card">
                <h3>{p.name}</h3>
                <p>{p.summary}</p>
                <ul className="bullets">
                  {p.bullets.map((b) => (
                    <li key={b}>{b}</li>
                  ))}
                </ul>
              </article>
            ))}
          </div>
        )}
      </section>

      <section id="education">
        <h2>Education</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : (
          <div className="edu-grid">
            {(data?.education || []).map((ed) => (
              <article key={ed.id} className="card">
                <h3>{ed.degree}</h3>
                <p className="muted">{ed.institution}</p>
                <p>
                  {ed.location} · {formatRange(ed.start_date, ed.end_date)}
                </p>
              </article>
            ))}
          </div>
        )}
      </section>

      <section id="contact">
        <h2>Contact</h2>
        {loading ? (
          <p className="muted">Loading…</p>
        ) : (
          <div className="contact-wrap">
            <p>
              REST <code>POST /api/v1/contact</code> persists to MySQL, then a worker pool publishes{" "}
              <code>contact.created</code> to RabbitMQ.
            </p>
            <form onSubmit={onSubmit} className="card form">
              <label>
                Name
                <input
                  required
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                />
              </label>
              <label>
                Email
                <input
                  required
                  type="email"
                  value={form.email}
                  onChange={(e) => setForm({ ...form, email: e.target.value })}
                />
              </label>
              <label>
                Subject
                <input
                  required
                  value={form.subject}
                  onChange={(e) => setForm({ ...form, subject: e.target.value })}
                />
              </label>
              <label>
                Message
                <textarea
                  required
                  rows="5"
                  value={form.message}
                  onChange={(e) => setForm({ ...form, message: e.target.value })}
                />
              </label>
              <button type="submit">Send message</button>
              {status && <p className="status">{status}</p>}
            </form>
          </div>
        )}
      </section>

      <footer>
        <span>Stack: Go · Gin · MySQL · Redis · RabbitMQ · React</span>
        <span>Cached reads on GET /api/v1/portfolio</span>
      </footer>
    </div>
  );
}

function Stat({ value, label }) {
  return (
    <div className="stat">
      <strong>{value}</strong>
      <span>{label}</span>
    </div>
  );
}
