import {type Badge, type Category} from '../models'

type Props = {
  categories: Category[]
  badgeOptions: Badge[]
  category: string
  status: string
  badges: string[]
  badgeInput: string
  onCategoryChange: (value: string) => void
  onStatusChange: (value: string) => void
  onBadgesChange: (badges: string[]) => void
  onBadgeInputChange: (value: string) => void
}

export function ContentFilters(props: Props) {
  const addBadge = () => {
    const badge = props.badgeInput.trim().toLowerCase()
    if (badge && !props.badges.includes(badge)) props.onBadgesChange([...props.badges, badge])
    props.onBadgeInputChange('')
  }
  return <div className="filters" aria-label="Content filters">
    <label className="status-filter">Status
      <select value={props.status} onChange={event => props.onStatusChange(event.target.value)}>
        <option value="">All statuses</option><option value="published">Published</option>
        <option value="draft">Draft</option><option value="review">Review</option><option value="archived">Archived</option>
      </select>
    </label>
    <label className="category-filter">Category
      <select value={props.category} onChange={event => props.onCategoryChange(event.target.value)}>
        <option value="">All categories</option>
        {props.categories.map(item => <option value={item.slug} key={item.slug}>{item.name}</option>)}
      </select>
    </label>
    <label className="badge-filter">Badges
      <div className="badge-input">
        {props.badges.map(badge => <span className="filter-pill" key={badge}>{badge}
          <button type="button" aria-label={`Remove ${badge}`}
            onClick={() => props.onBadgesChange(props.badges.filter(value => value !== badge))}>×</button></span>)}
        <input list="badge-options" value={props.badgeInput} placeholder="Type a badge and press Enter"
          onChange={event => props.onBadgeInputChange(event.target.value)}
          onKeyDown={event => { if (event.key === 'Enter') { event.preventDefault(); addBadge() } }}/>
        <datalist id="badge-options">{props.badgeOptions.filter(item => !props.badges.includes(item.name)).map(item =>
          <option value={item.name} key={item.id}>{item.name}</option>)}</datalist>
      </div>
    </label>
  </div>
}
