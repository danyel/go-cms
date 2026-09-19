-- Seed data for the PostgreSQL content table.
-- Safe to run repeatedly: rows are upserted by slug.
--
-- Run with:
--   psql "$POSTGRES_URL" -f scripts/seed_content.sql

INSERT INTO categories (slug, name) VALUES
    ('engineering', 'Engineering'),
    ('design', 'Design'),
    ('product', 'Product'),
    ('editorial', 'Editorial'),
    ('community', 'Community')
ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO content (
    slug, title, summary, body, status, published, created_at, updated_at
) VALUES
(
    'welcome-to-the-cms',
    'Welcome to the CMS',
    'A short introduction to the editorial workspace and the ideas behind this publication.',
    $$Welcome to the CMS journal.

This space is where the team shares product updates, practical guides, and notes from the work behind the scenes. Each story is written to be useful on its own, while also giving readers a clearer picture of how the product is built.

Start with the latest articles, save the pieces that are useful to your team, and send us feedback when something could be clearer.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '24 days', CURRENT_TIMESTAMP - INTERVAL '20 days'
),
(
    'designing-for-clarity',
    'Designing for Clarity',
    'How a small set of deliberate design decisions can make a complex workflow feel calm and understandable.',
    $$Clarity is not the same as simplicity. A good interface can expose a powerful workflow without asking people to understand every implementation detail.

We begin by naming the user''s goal, then remove everything that competes with it. Labels should describe outcomes, actions should have predictable results, and important feedback should appear close to the action that caused it.

The best design review question is often simple: what would a first-time user expect to happen next?$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '23 days', CURRENT_TIMESTAMP - INTERVAL '18 days'
),
(
    'a-practical-guide-to-editorial-workflows',
    'A Practical Guide to Editorial Workflows',
    'A repeatable process for taking an idea from a rough note to a polished, published story.',
    $$A dependable editorial workflow gives ideas room to develop without allowing unfinished work to slip into production.

Our process has four stages: capture, shape, review, and publish. Capture is intentionally lightweight. Shape gives the draft a clear audience and purpose. Review checks accuracy, tone, and accessibility. Publish is the final technical and editorial check.

The workflow works because every stage has a clear owner and a clear definition of done.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '22 days', CURRENT_TIMESTAMP - INTERVAL '16 days'
),
(
    'building-a-better-content-calendar',
    'Building a Better Content Calendar',
    'A lightweight planning method that balances consistency with enough flexibility for timely ideas.',
    $$A content calendar should create confidence, not bureaucracy.

We plan around themes rather than filling every date in advance. Each month has a primary theme, two supporting topics, and a small reserve for timely announcements. This gives the team a reliable rhythm while leaving space for unexpected opportunities.

The calendar is most useful when it shows the reason for each piece, not only its due date.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '21 days', CURRENT_TIMESTAMP - INTERVAL '15 days'
),
(
    'the-case-for-small-releases',
    'The Case for Small Releases',
    'Why shipping focused improvements regularly creates better feedback and lower operational risk.',
    $$Small releases make learning cheaper.

When a change has a narrow scope, it is easier to test, explain, monitor, and roll back. The team gets useful feedback while the context is still fresh, and users see a steady stream of improvements rather than waiting for a single large launch.

Small does not mean insignificant. A series of focused changes can transform an experience when each one solves a real problem.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '20 days', CURRENT_TIMESTAMP - INTERVAL '14 days'
),
(
    'how-we-write-useful-documentation',
    'How We Write Useful Documentation',
    'Documentation is most valuable when it helps someone complete a task at the moment they need it.',
    $$Useful documentation starts with the reader''s question, not the system''s architecture.

We lead with the expected result, show the smallest complete example, and explain the assumptions that could otherwise cause confusion. Details belong after the main path, where they can help without blocking progress.

Documentation is part of the product. When it is accurate, direct, and easy to scan, it reduces support work and helps users become independent faster.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '19 days', CURRENT_TIMESTAMP - INTERVAL '13 days'
),
(
    'accessibility-is-a-quality-practice',
    'Accessibility Is a Quality Practice',
    'Accessibility checks improve the experience for everyone, not only for users who identify with a specific access need.',
    $$Accessibility is a practical way to test whether an interface communicates clearly.

Keyboard navigation reveals hidden interaction assumptions. Good contrast helps in bright rooms and on small screens. Descriptive labels help screen-reader users and also make automated tests more meaningful.

Treating accessibility as a quality practice means including it in design, implementation, and review instead of leaving it for the final checklist.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '18 days', CURRENT_TIMESTAMP - INTERVAL '12 days'
),
(
    'measuring-what-matters',
    'Measuring What Matters',
    'A guide to choosing a small set of product signals that lead to better decisions.',
    $$Metrics are useful when they change what the team does.

We prefer a small collection of signals with clear owners: successful task completion, time to value, reliable publishing, and meaningful engagement. Every metric should have a question behind it and a decision that it can inform.

A dashboard full of numbers is not automatically useful. The goal is shared understanding, not decorative precision.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '17 days', CURRENT_TIMESTAMP - INTERVAL '11 days'
),
(
    'notes-from-a-content-migration',
    'Notes from a Content Migration',
    'Lessons learned while moving a growing library into a structured content model.',
    $$A migration is an opportunity to clarify the shape of the content, not only to move records from one place to another.

We started by identifying the fields that readers actually depend on. Then we separated durable content from presentation details and created a small set of validation rules. The result was less data than we expected, but much more confidence in what remained.

The most important migration tool was a clear decision log. It made trade-offs visible and prevented the team from reopening settled questions.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '16 days', CURRENT_TIMESTAMP - INTERVAL '10 days'
),
(
    'making-search-feel-natural',
    'Making Search Feel Natural',
    'Search works best when it understands the words people use, not only the labels in a database.',
    $$People search with questions, fragments, and remembered phrases. They rarely use the exact vocabulary that an internal system uses.

Useful search begins with good titles and summaries, then adds ranking that rewards relevance without hiding recent work. Empty results should offer a next step, and near matches should be understandable rather than mysterious.

Search is a conversation with the content library. Every result teaches users what the system considers important.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '15 days', CURRENT_TIMESTAMP - INTERVAL '9 days'
),
(
    'a-field-guide-to-content-reviews',
    'A Field Guide to Content Reviews',
    'A concise review checklist for accuracy, usefulness, voice, and technical readiness.',
    $$A strong review is more than proofreading.

First confirm that the piece answers a real reader question. Then verify names, dates, links, examples, and claims. Check that the heading structure is logical and that the summary makes sense outside the article. Finally, read it once as a person who knows nothing about the project.

Reviewers should improve the work while protecting the author''s intent. Specific suggestions are more helpful than broad requests to make something better.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '14 days', CURRENT_TIMESTAMP - INTERVAL '8 days'
),
(
    'reliability-behind-the-scenes',
    'Reliability Behind the Scenes',
    'The habits and safeguards that help a small team operate dependable services.',
    $$Reliability is built from ordinary habits practiced consistently.

We keep operational changes small, make failures visible, and write down the recovery path before an incident happens. Health checks tell us whether a service is alive; useful logs help us understand whether it is doing the right work.

The goal is not to eliminate every failure. It is to make failures understandable, contained, and recoverable.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '13 days', CURRENT_TIMESTAMP - INTERVAL '7 days'
),
(
    'the-anatomy-of-a-good-release-note',
    'The Anatomy of a Good Release Note',
    'Release notes should help readers understand what changed and why it matters to them.',
    $$A good release note answers three questions quickly: what changed, who benefits, and what should the reader do now?

We avoid internal ticket numbers and implementation language unless they clarify the result. A short example is often more useful than a list of changed files. When a change requires action, the action should be explicit and easy to find.

Release notes are part of the product''s memory. Written well, they help current users and future teammates understand how the experience evolved.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '12 days', CURRENT_TIMESTAMP - INTERVAL '6 days'
),
(
    'working-with-feedback',
    'Working with Feedback',
    'A practical approach to turning scattered comments into decisions and improvements.',
    $$Feedback is raw material, not a verdict.

We group comments by the problem they describe, look for repeated evidence, and separate requests from the underlying need. Some feedback should lead to an immediate fix; some should become a research question; some should be acknowledged and intentionally deferred.

The most useful response to feedback is a visible decision. People do not need every request to be accepted, but they do need to know that it was understood.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '11 days', CURRENT_TIMESTAMP - INTERVAL '5 days'
),
(
    'the-value-of-a-calm-admin-interface',
    'The Value of a Calm Admin Interface',
    'Administrative tools should reduce cognitive load so people can focus on making good decisions.',
    $$An admin interface is often used under time pressure, which makes calm design especially important.

Clear status labels, predictable controls, and safe defaults help editors move quickly without guessing. Destructive actions need friction, while routine publishing actions should feel direct and reversible.

Good administration is invisible when it works. The interface gives people confidence without demanding their attention.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '4 days'
),
(
    'when-to-use-a-draft',
    'When to Use a Draft',
    'Drafts are a tool for thinking in public without publishing unfinished decisions.',
    $$A draft is useful whenever the content or the decision still needs context.

Use a draft to collect research, test an outline, ask for review, or preserve an idea that is not ready for an audience. A draft should still have a clear next step; otherwise it becomes a storage drawer.

The best draft is easy for another person to pick up. Give it a useful title, explain what is unresolved, and note what kind of feedback would help.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '9 days', CURRENT_TIMESTAMP - INTERVAL '3 days'
),
(
    'our-approach-to-technical-debt',
    'Our Approach to Technical Debt',
    'Technical debt is easier to manage when it is connected to a user impact or an operational risk.',
    $$Not every shortcut needs immediate repayment. The important question is whether the shortcut is making future work slower, riskier, or harder to understand.

We record debt near the code or workflow it affects, describe the cost in concrete terms, and revisit it when related work is already planned. This keeps maintenance connected to product decisions rather than turning it into an abstract list.

Small, regular payments are usually healthier than a dramatic cleanup project.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '2 days'
),
(
    'a-better-way-to-plan-content',
    'A Better Way to Plan Content',
    'Planning becomes easier when every proposed story has a reader, a purpose, and a natural next step.',
    $$Content planning is not a contest to produce the most ideas.

For each idea, write down who it helps, what question it answers, and what the reader can do afterward. This makes promising ideas easier to prioritize and exposes ideas that are interesting but not yet useful.

The strongest plans leave room for discovery. A clear purpose is more valuable than a rigid publishing schedule.$$,
    'review', FALSE, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '1 day'
),
(
    'meet-the-editorial-team',
    'Meet the Editorial Team',
    'A behind-the-scenes look at the people who shape the CMS journal.',
    $$The editorial team brings together product, design, engineering, and customer perspectives.

That mix helps us ask better questions. Writers bring clarity, subject experts bring depth, and reviewers make sure the result is accurate and useful. The team is intentionally small so that conversations stay close to the work.

In the coming weeks we will share more of the practices that help this group collaborate.$$,
    'draft', FALSE, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '12 hours'
),
(
    'what-we-learned-from-our-first-workshop',
    'What We Learned from Our First Workshop',
    'A reflection on the exercises that helped a mixed team agree on the most important user problems.',
    $$Our first workshop was designed to create alignment without pretending that every question had an immediate answer.

The most productive exercise was mapping the moments where users hesitate. Those moments gave the team a shared language for discussing confusing workflows and showed where a small improvement could have an outsized effect.

The next step is to turn those observations into experiments that can be tested with real users.$$,
    'draft', FALSE, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '8 hours'
),
(
    'a-small-experiment-in-onboarding',
    'A Small Experiment in Onboarding',
    'We are testing a shorter first-run experience focused on one successful outcome.',
    $$This is an early experiment, not a final design.

The current version asks new users to complete one meaningful task before introducing secondary features. We expect this to make the first session more focused and give us a clearer signal about where people need help.

We will compare completion rates and qualitative feedback before deciding whether to keep the change.$$,
    'review', FALSE, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '4 hours'
),
(
    'questions-we-are-exploring',
    'Questions We Are Exploring',
    'A living list of questions that will guide upcoming product and editorial work.',
    $$We are exploring how teams decide which content deserves an update, how readers move from one useful story to the next, and how editors can collaborate without losing ownership.

These questions are intentionally open. They will become more specific as we talk to users, observe workflows, and test small improvements.

If your team has a different perspective, we would like to hear what you have learned.$$,
    'draft', FALSE, CURRENT_TIMESTAMP - INTERVAL '18 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'
),
(
    'content-modeling-principles',
    'Content Modeling Principles',
    'The principles we use to keep content structured, portable, and easy to evolve.',
    $$A content model should reflect the meaning of content before it reflects a particular screen.

We prefer explicit fields over hidden conventions, stable identifiers over presentation-dependent names, and validation that explains what is wrong. A model should support today''s workflows while leaving room for tomorrow''s channels.

The best model is not the most abstract one. It is the one that helps people create and maintain useful content with confidence.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '7 days', CURRENT_TIMESTAMP - INTERVAL '1 day'
),
(
    'why-plain-language-wins',
    'Why Plain Language Wins',
    'Clear writing helps people make decisions faster and makes products feel more trustworthy.',
    $$Plain language is not about removing nuance. It is about putting the nuance where readers can use it.

Use familiar words, prefer active sentences, and explain specialized terms when they matter. Short paragraphs make ideas easier to navigate, while concrete examples help readers connect advice to their own work.

When writing is clear, readers spend their energy on the decision instead of decoding the sentence.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '6 days', CURRENT_TIMESTAMP - INTERVAL '20 hours'
),
(
    'a-note-on-data-retention',
    'A Note on Data Retention',
    'How we think about keeping only the data that supports a clear product or operational purpose.',
    $$Data should have a reason to exist and a responsible owner.

We review retention from the perspective of usefulness, privacy, and operational cost. When data no longer serves a purpose, deleting it is often safer than preserving it just in case.

Good retention practices are part of good product design because they make expectations clearer for everyone involved.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '16 hours'
),
(
    'the-next-quarter',
    'Looking Ahead to the Next Quarter',
    'The themes that will shape our next set of product, platform, and editorial improvements.',
    $$The next quarter is about making the core experience more dependable and more useful.

We will focus on faster editorial workflows, clearer content discovery, stronger operational feedback, and a more consistent foundation for future features. We will measure progress by the quality of completed work, not by the number of initiatives started.

There is a lot to do, so focus will matter as much as ambition.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '8 hours'
),
(
    'community-notes-june',
    'Community Notes: June',
    'A monthly collection of questions, ideas, and observations shared by the CMS community.',
    $$This month the community asked for better preview tools, clearer publishing status, and more ways to reuse content across channels.

Several teams also shared small workflow improvements that they built locally. Those examples remind us that the best ideas often begin as a practical response to a real constraint.

We are collecting these notes to make useful patterns easier to find and discuss.$$,
    'published', TRUE, CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '6 hours'
),
(
    'release-0-3-preview',
    'Release 0.3 Preview',
    'A preview of the next CMS iteration, including improved content editing and migration tooling.',
    $$The next release will make content editing more approachable and the local development workflow more predictable.

Planned improvements include clearer editor permissions, stronger migration checks, better empty states, and a more useful content preview. These changes are still being tested and may change before release.

We will publish final notes when the release is ready for general use.$$,
    'review', FALSE, CURRENT_TIMESTAMP - INTERVAL '12 hours', CURRENT_TIMESTAMP - INTERVAL '1 hour'
)
ON CONFLICT (slug) DO UPDATE SET
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    body = EXCLUDED.body,
    status = EXCLUDED.status,
    published = EXCLUDED.published,
    updated_at = EXCLUDED.updated_at;

UPDATE content SET category_id = categories.id
FROM categories
WHERE categories.slug = CASE
    WHEN content.slug IN ('designing-for-clarity', 'accessibility-is-a-quality-practice', 'the-value-of-a-calm-admin-interface') THEN 'design'
    WHEN content.slug IN ('the-case-for-small-releases', 'measuring-what-matters', 'making-search-feel-natural', 'a-small-experiment-in-onboarding', 'the-next-quarter') THEN 'product'
    WHEN content.slug IN ('a-practical-guide-to-editorial-workflows', 'building-a-better-content-calendar', 'how-we-write-useful-documentation', 'a-field-guide-to-content-reviews', 'when-to-use-a-draft', 'a-better-way-to-plan-content', 'why-plain-language-wins') THEN 'editorial'
    WHEN content.slug IN ('meet-the-editorial-team', 'working-with-feedback', 'community-notes-june', 'what-we-learned-from-our-first-workshop', 'questions-we-are-exploring') THEN 'community'
    ELSE 'engineering'
END;

INSERT INTO badges (name) VALUES
    ('go'), ('postgresql'), ('docker'), ('react'), ('typescript'),
    ('java'), ('spring-boot'), ('maven'), ('accessibility'),
    ('content-strategy'), ('product-management'), ('documentation'),
    ('reliability'), ('search'), ('community')
ON CONFLICT (name) DO NOTHING;

INSERT INTO content_badges (content_id, badge_id)
SELECT c.id, b.id
FROM (VALUES
    ('welcome-to-the-cms', 'content-strategy'),
    ('welcome-to-the-cms', 'community'),
    ('designing-for-clarity', 'accessibility'),
    ('designing-for-clarity', 'content-strategy'),
    ('a-practical-guide-to-editorial-workflows', 'content-strategy'),
    ('building-a-better-content-calendar', 'content-strategy'),
    ('the-case-for-small-releases', 'product-management'),
    ('how-we-write-useful-documentation', 'documentation'),
    ('accessibility-is-a-quality-practice', 'accessibility'),
    ('measuring-what-matters', 'product-management'),
    ('notes-from-a-content-migration', 'postgresql'),
    ('notes-from-a-content-migration', 'docker'),
    ('making-search-feel-natural', 'search'),
    ('a-field-guide-to-content-reviews', 'documentation'),
    ('reliability-behind-the-scenes', 'go'),
    ('reliability-behind-the-scenes', 'docker'),
    ('the-anatomy-of-a-good-release-note', 'documentation'),
    ('working-with-feedback', 'community'),
    ('the-value-of-a-calm-admin-interface', 'design'),
    ('when-to-use-a-draft', 'content-strategy'),
    ('our-approach-to-technical-debt', 'go'),
    ('our-approach-to-technical-debt', 'postgresql'),
    ('a-better-way-to-plan-content', 'content-strategy'),
    ('meet-the-editorial-team', 'community'),
    ('what-we-learned-from-our-first-workshop', 'community'),
    ('a-small-experiment-in-onboarding', 'react'),
    ('a-small-experiment-in-onboarding', 'typescript'),
    ('questions-we-are-exploring', 'product-management'),
    ('content-modeling-principles', 'postgresql'),
    ('content-modeling-principles', 'content-strategy'),
    ('why-plain-language-wins', 'documentation'),
    ('a-note-on-data-retention', 'postgresql'),
    ('the-next-quarter', 'product-management'),
    ('community-notes-june', 'community'),
    ('release-0-3-preview', 'react'),
    ('release-0-3-preview', 'go')
) AS seed(slug, badge)
JOIN content c ON c.slug = seed.slug
JOIN badges b ON b.name = seed.badge
ON CONFLICT DO NOTHING;
