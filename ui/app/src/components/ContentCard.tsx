import {Link} from 'react-router-dom'
import {type Content} from '../models'

export function ContentCard({item}: {item: Content}) {
  return <article className="content-card">
    <span className="status">{item.status}</span>
    <h2><Link to={`/content/${item.id}`}>{item.title}</Link></h2>
    <p>{item.summary}</p>
    <p className="content-meta">by {item.author ?? 'unknown'}
      {item.publishedAt && <> · published {new Date(item.publishedAt).toLocaleDateString()}</>}</p>
    {item.badges && item.badges.length > 0 &&
      <div className="badge-list">{item.badges.map(badge => <span className="badge" key={badge}>{badge}</span>)}</div>}
    <Link className="read-link" to={`/content/${item.id}`}>Read more <span aria-hidden="true">→</span></Link>
  </article>
}
