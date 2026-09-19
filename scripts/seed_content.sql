-- Seed data for the PostgreSQL content table.
-- Safe to run repeatedly: rows are upserted by slug.
--
-- Run with:
--   psql "$POSTGRES_URL" -f scripts/seed_content.sql

INSERT INTO categories (slug, name) VALUES
    ('linux', 'Linux'),
    ('bash', 'Bash & Shell'),
    ('go', 'Go'),
    ('java', 'Java'),
    ('devops', 'DevOps')
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

-- Re-theme the repeatable examples around Linux, shell tooling, Go, Java, and DevOps.
INSERT INTO badges (name) VALUES
    ('linux'), ('bash'), ('shell'), ('terminal'), ('systemd'), ('grep'),
    ('awk'), ('sed'), ('ssh'), ('git'), ('journalctl'), ('go'), ('java'), ('spring-boot'),
    ('maven'), ('postgresql'), ('docker'), ('kubernetes'), ('devops')
ON CONFLICT (name) DO NOTHING;

UPDATE content AS c
SET title = seed.title, summary = seed.summary, body = seed.body
FROM (VALUES
    ('welcome-to-the-cms', 'Welcome to Urpi''s backlog', 'A terminal-shaped notebook for Linux, Bash, Go, Java, and useful experiments.', $$Welcome to Urpi''s backlog.

This is a working notebook for the tools that make a developer''s day better: a clean Linux shell, a small Bash script, a reliable Go service, and a Java application that is easier to operate.

The best notes are practical. They show the command, explain the why, and leave enough context for future-you to use them again.$$),
    ('designing-for-clarity', 'Designing a Calm Terminal Workflow', 'Small conventions that make command-line work easier to read, repeat, and recover.', $$A good terminal workflow is easy to scan at a glance.

Use clear prompts, predictable directories, short commands, and output that explains what happened. Keep destructive operations explicit and make the safe path the easy path.

The shell is not only a place to execute commands. It is an interface worth designing.$$),
    ('a-practical-guide-to-editorial-workflows', 'A Practical Bash Script Workflow', 'How to write shell scripts that are safe to run, easy to debug, and friendly to the next maintainer.', $$A useful Bash script starts with strict mode, meaningful names, and a clear exit path.

Use set -euo pipefail when appropriate, quote variables, validate inputs, and print progress around operations that may take time. Prefer small functions over one long pipeline when the script is important enough to maintain.

Scripts become tools when another person can understand them without reverse engineering the author''s terminal session.$$),
    ('building-a-better-content-calendar', 'Building a Personal Linux Command Log', 'A simple way to turn repeated terminal discoveries into searchable notes.', $$Whenever a command solves a problem, write down the context with it.

Record the operating system, the input, the useful output, and the reason the command worked. A note such as journalctl -u service --since today is much more valuable when it says which service was failing and what the logs revealed.

Over time, the command log becomes a small, personal operations manual.$$),
    ('the-case-for-small-releases', 'The Case for Small CLI Releases', 'Why focused command-line improvements are easier to test, ship, and roll back.', $$A small release is easier to inspect with git diff, test from a clean shell, and undo when it does not behave as expected.

This is especially useful for developer tools. A narrow change to a flag or output format gives users a clear upgrade path and gives maintainers a focused set of failure modes.

Small releases keep the feedback loop close to the command that changed.$$),
    ('how-we-write-useful-documentation', 'How to Document a Bash Command', 'A command is not documented until someone knows what it changes and how to undo it.', $$Good shell documentation shows a complete example, explains expected output, and calls out permissions or paths that matter.

Always distinguish a read-only command from one that modifies files or services. Include a dry-run option when possible and show the recovery command beside the change.

The goal is not more prose. It is a safer next terminal session.$$),
    ('accessibility-is-a-quality-practice', 'Readable Output Is a Quality Practice', 'Terminal output should work in bright rooms, small windows, logs, and automated checks.', $$Readable output uses stable labels, useful exit codes, and enough spacing to separate one result from the next.

Do not rely on color alone. A warning should still be clear when output is redirected to a file or viewed through SSH. Good command-line output helps both humans and scripts.

Accessibility is another word for making the right information available in the conditions people actually work in.$$),
    ('measuring-what-matters', 'Measuring a Service from the Shell', 'A few dependable commands can tell you more than a noisy dashboard.', $$Start with the signals that answer an operational question: is the process running, is the port listening, are requests succeeding, and is the disk filling up?

Commands such as systemctl status, ss, curl, df, and journalctl are powerful because they connect directly to the machine''s state.

Measure what helps you choose the next command, not what merely looks impressive.$$),
    ('notes-from-a-content-migration', 'Notes from a PostgreSQL Migration', 'Lessons from moving local content data into PostgreSQL with repeatable migrations.', $$A database migration should be boring to run and easy to inspect.

Keep schema changes in versioned SQL, make seed data idempotent, and verify the result with psql before the application starts depending on it. A migration that works only on one laptop is not finished.

The terminal is an excellent place to make database state visible.$$),
    ('making-search-feel-natural', 'Searching Code with grep, rg, and git', 'A practical tour of fast text searches for a busy Linux repository.', $$Start broad with rg, narrow by file type, and include line numbers when the result will become a follow-up task.

git grep is useful when you want the tracked view of a repository. grep -R remains available almost everywhere, which makes it a dependable fallback on minimal systems.

The best search command shortens the distance between a question and the file that answers it.$$),
    ('a-field-guide-to-content-reviews', 'A Field Guide to Shell Script Reviews', 'A compact checklist for reviewing Bash before it reaches production.', $$Check quoting, glob behavior, unset variables, exit codes, temporary files, permissions, and cleanup paths.

Then read the script as a sequence of commands run by a new shell on a machine with a different environment. Reviewers should ask what happens when the network is unavailable, a file is missing, or a command returns no rows.

The best review catches an unsafe assumption before it becomes an incident.$$),
    ('reliability-behind-the-scenes', 'Reliability with systemd and journalctl', 'A practical pairing for running and diagnosing Linux services.', $$systemd describes how a service should run; journalctl shows what it actually did.

Use systemctl status for a quick state check, systemctl cat to inspect the loaded unit, and journalctl -u service --since today to follow the evidence. Keep restart policies intentional and make health checks observable.

Reliable services leave a useful trail when something goes wrong.$$),
    ('the-anatomy-of-a-good-release-note', 'The Anatomy of a Good CLI Release', 'Release notes for developer tools should show the changed command and its effect.', $$A useful CLI release note includes the old behavior, the new command, and one realistic terminal example.

Mention changed defaults, removed flags, migration steps, and compatibility concerns. Users should be able to decide whether to upgrade without searching through implementation details.

Good release notes make a tool feel stable even as it evolves.$$),
    ('working-with-feedback', 'Working with Terminal Feedback', 'How command output, bug reports, and shell history reveal what a tool should improve.', $$A failed command is feedback with context attached.

Capture the exact command, environment, exit status, and the smallest useful piece of output. Repeated confusion around the same flag is often a documentation problem; repeated slow commands may be a design problem.

Treat reports as clues about the workflow, not only as isolated defects.$$),
    ('the-value-of-a-calm-admin-interface', 'The Value of a Calm Operations Console', 'A good operations screen should feel like a well-organized terminal session.', $$Operations work needs clear status, safe defaults, and an obvious path back from a mistake.

Show what is running, what changed, and which command or action will happen next. Destructive operations should make scope visible before they ask for confirmation.

The calmest interface is the one that helps a tired operator make the right decision.$$),
    ('when-to-use-a-draft', 'When to Keep a Shell Note as a Draft', 'Some commands need a little more testing before they become part of the permanent runbook.', $$Keep a command note in draft while paths, permissions, and failure behavior are still uncertain.

Run it from a clean directory, test the negative path, and write down the rollback before calling it ready. A draft is useful when it preserves the question as well as the current answer.$$),
    ('our-approach-to-technical-debt', 'Our Approach to Shell and Build Debt', 'Technical debt in scripts and build files compounds quickly when nobody owns the cleanup.', $$A clever one-liner can become expensive when it is copied into five deployment scripts.

Prefer named functions, stable interfaces, and explicit dependencies. In Go and Java projects, keep build commands discoverable through Make, Maven, or a documented task runner.

Pay down debt when you are already working near the command that created it.$$),
    ('a-better-way-to-plan-content', 'A Better Way to Plan a Linux Lab', 'Plan experiments around one command, one hypothesis, and one observable result.', $$A useful Linux lab has a question such as “what happens when this unit restarts?” and a small set of commands that can answer it.

Write the expected result before running the experiment. Save the actual output, explain the difference, and turn the conclusion into a reusable note.

Focused experiments produce better operational knowledge than a long list of disconnected commands.$$),
    ('meet-the-editorial-team', 'Meet the Tools in the Terminal', 'The daily toolkit: Bash for glue, Go for services, Java for platforms, and Linux underneath.', $$Different tools earn their place by solving different problems.

Bash connects existing commands. Go makes small services and utilities easy to ship. Java and Spring Boot support long-lived applications with strong conventions. Linux gives all of them a dependable home.

The interesting work happens at the boundaries between those tools.$$),
    ('what-we-learned-from-our-first-workshop', 'What We Learned from a Shell Workshop', 'A few exercises that helped developers become more confident at the command line.', $$The most useful workshop exercise was troubleshooting a broken service from symptoms alone.

Participants used systemctl, journalctl, ss, curl, and grep to form a hypothesis, test it, and explain the fix. The lesson was less about memorizing commands and more about building a reliable investigation loop.

Confidence grows when every command has a question behind it.$$),
    ('a-small-experiment-in-onboarding', 'A Small Experiment with Bash Completion', 'We are testing whether better completion makes unfamiliar commands easier to discover.', $$Completion can turn a blank prompt into a gentle guide.

The experiment adds examples for common flags, paths, and service names while keeping the underlying command unchanged. We will compare successful first attempts and the number of help-page lookups.

Good completion should teach without getting in the way.$$),
    ('questions-we-are-exploring', 'Questions We Are Exploring', 'Open questions about Linux workflows, build tools, and the boundary between scripts and software.', $$When should a Bash script become a Go command? When does a Maven plugin deserve a dedicated build step? Which operational knowledge belongs in a runbook and which belongs in automation?

These questions are intentionally open. Experiments in the terminal will give us better answers than abstract rules.$$),
    ('content-modeling-principles', 'Content Modeling for Runbooks', 'A runbook needs commands, assumptions, expected output, and recovery steps—not just a paragraph of advice.', $$A useful runbook entry describes the situation, the command to run, the expected evidence, and the safe recovery path.

Keep environment-specific values visible, separate read-only checks from changes, and include links to the service or migration that the command supports.

Structured operational notes are easier to search when the next incident arrives.$$),
    ('why-plain-language-wins', 'Why Plain Language Wins in Man Pages', 'Clear command descriptions help people move from a prompt to a correct action.', $$A man page should tell readers what a command does before listing every option.

Use examples that resemble real work, explain dangerous flags near the command that uses them, and avoid hiding the important behavior behind jargon.

Plain language is a reliability feature for people working under pressure.$$),
    ('a-note-on-data-retention', 'A Note on Logs and Data Retention', 'Keep logs long enough to investigate, but not so long that useful evidence disappears in noise.', $$Log retention is an operational decision.

Choose a period that supports debugging and compliance, rotate files predictably, and make sure disk usage is observable. journalctl and logrotate are only useful when their limits are understood.

The right amount of history is the amount that helps answer the next incident question.$$),
    ('the-next-quarter', 'Looking Ahead from the Terminal', 'The next stretch of work focuses on sharper Linux workflows, dependable Go services, and practical Java notes.', $$The backlog ahead includes systemd troubleshooting, Bash patterns, Go HTTP services, Spring Boot operations, Maven build hygiene, and better PostgreSQL tooling.

The goal is not to collect technologies. It is to understand the small commands and decisions that make them pleasant to use.$$),
    ('community-notes-june', 'Community Notes from the Linux Shell', 'A collection of commands and habits shared by people who spend their days in terminals.', $$This month''s notes include using ssh config aliases, pairing find with xargs carefully, checking ports with ss, and using git worktree for parallel fixes.

The common thread is respect for the next person at the prompt: make the command visible, explain the edge case, and leave the machine in a known state.$$),
    ('release-0-3-preview', 'Release 0.3 Preview: More Terminal, Less Guessing', 'A preview of improvements for Linux-first development workflows.', $$The next release will add clearer command output, safer local migrations, better Go service diagnostics, and more useful notes for Java and Spring Boot projects.

These changes are still being tested. The target is simple: fewer guesses between opening a terminal and understanding what the application is doing.$$)
) AS seed(slug, title, summary, body) WHERE c.slug = seed.slug;

UPDATE content SET category_id = categories.id
FROM categories
WHERE categories.slug = CASE
    WHEN content.slug IN ('a-practical-guide-to-editorial-workflows', 'building-a-better-content-calendar', 'how-we-write-useful-documentation', 'accessibility-is-a-quality-practice', 'making-search-feel-natural', 'a-field-guide-to-content-reviews', 'when-to-use-a-draft', 'a-better-way-to-plan-content', 'why-plain-language-wins', 'welcome-to-the-cms') THEN 'bash'
    WHEN content.slug IN ('the-case-for-small-releases', 'measuring-what-matters', 'reliability-behind-the-scenes', 'our-approach-to-technical-debt', 'meet-the-editorial-team', 'what-we-learned-from-our-first-workshop') THEN 'go'
    WHEN content.slug IN ('a-small-experiment-in-onboarding', 'questions-we-are-exploring') THEN 'java'
    WHEN content.slug IN ('notes-from-a-content-migration', 'the-anatomy-of-a-good-release-note', 'working-with-feedback', 'the-value-of-a-calm-admin-interface', 'content-modeling-principles', 'a-note-on-data-retention', 'the-next-quarter', 'community-notes-june', 'release-0-3-preview') THEN 'devops'
    ELSE 'linux'
END;

DELETE FROM content_badges;
INSERT INTO content_badges (content_id, badge_id)
SELECT c.id, b.id
FROM (VALUES
    ('welcome-to-the-cms', 'linux'), ('welcome-to-the-cms', 'terminal'), ('welcome-to-the-cms', 'bash'),
    ('designing-for-clarity', 'terminal'), ('designing-for-clarity', 'shell'),
    ('a-practical-guide-to-editorial-workflows', 'bash'), ('a-practical-guide-to-editorial-workflows', 'shell'),
    ('building-a-better-content-calendar', 'linux'), ('building-a-better-content-calendar', 'grep'),
    ('the-case-for-small-releases', 'go'), ('the-case-for-small-releases', 'git'),
    ('how-we-write-useful-documentation', 'bash'), ('how-we-write-useful-documentation', 'terminal'),
    ('accessibility-is-a-quality-practice', 'linux'), ('accessibility-is-a-quality-practice', 'shell'),
    ('measuring-what-matters', 'linux'), ('measuring-what-matters', 'systemd'),
    ('notes-from-a-content-migration', 'postgresql'), ('notes-from-a-content-migration', 'docker'),
    ('making-search-feel-natural', 'grep'), ('making-search-feel-natural', 'awk'),
    ('a-field-guide-to-content-reviews', 'bash'), ('a-field-guide-to-content-reviews', 'shell'),
    ('reliability-behind-the-scenes', 'systemd'), ('reliability-behind-the-scenes', 'linux'),
    ('the-anatomy-of-a-good-release-note', 'go'), ('the-anatomy-of-a-good-release-note', 'git'),
    ('working-with-feedback', 'terminal'), ('working-with-feedback', 'ssh'),
    ('the-value-of-a-calm-admin-interface', 'devops'), ('the-value-of-a-calm-admin-interface', 'linux'),
    ('when-to-use-a-draft', 'bash'), ('when-to-use-a-draft', 'git'),
    ('our-approach-to-technical-debt', 'go'), ('our-approach-to-technical-debt', 'java'),
    ('a-better-way-to-plan-content', 'linux'), ('a-better-way-to-plan-content', 'terminal'),
    ('meet-the-editorial-team', 'bash'), ('meet-the-editorial-team', 'go'), ('meet-the-editorial-team', 'java'),
    ('what-we-learned-from-our-first-workshop', 'systemd'), ('what-we-learned-from-our-first-workshop', 'journalctl'),
    ('a-small-experiment-in-onboarding', 'bash'), ('a-small-experiment-in-onboarding', 'terminal'),
    ('questions-we-are-exploring', 'go'), ('questions-we-are-exploring', 'spring-boot'), ('questions-we-are-exploring', 'maven'),
    ('content-modeling-principles', 'devops'), ('content-modeling-principles', 'postgresql'),
    ('why-plain-language-wins', 'linux'), ('why-plain-language-wins', 'terminal'),
    ('a-note-on-data-retention', 'systemd'), ('a-note-on-data-retention', 'devops'),
    ('the-next-quarter', 'linux'), ('the-next-quarter', 'go'), ('the-next-quarter', 'java'),
    ('community-notes-june', 'ssh'), ('community-notes-june', 'git'), ('community-notes-june', 'shell'),
    ('release-0-3-preview', 'go'), ('release-0-3-preview', 'spring-boot'), ('release-0-3-preview', 'maven')
) AS seed(slug, badge)
JOIN content c ON c.slug = seed.slug
JOIN badges b ON b.name = seed.badge;
