package service

import (
	"fmt"
	"strings"
	"time"
)

// categoryNames maps the seeded category slugs to their display names. The
// filter dropdown is built from the categories actually present in the seed and
// keeps this order.
var categoryNames = map[string]string{
	"linux":    "Linux",
	"bash":     "Bash & Shell",
	"go":       "Go",
	"java":     "Java",
	"devops":   "DevOps",
	"tooling":  "Tooling",
	"design":   "Design & Editorial",
	"workflow": "Workflow",
}

// daysAgo returns a timestamp relative to startup so the demo always shows a
// fresh-looking timeline.
func daysAgo(days int) time.Time {
	return time.Now().AddDate(0, 0, -days).Truncate(time.Second)
}

// SeedContent returns the demo backlog loaded into the in-memory store.
func SeedContent() []Content {
	items := make([]Content, 0, 560)
	for i, item := range seedContent {
		item.Author = demoAuthors[i%len(demoAuthors)]
		if item.Status == "published" {
			published := item.CreatedAt.Add(24 * time.Hour)
			item.PublishedAt = &published
		}
		items = append(items, item)
	}
	for i := len(items); i < 560; i++ {
		items = append(items, generatedContent(i))
	}
	return items
}

var demoAuthors = []string{"urpi", "tux", "gopher"}

func generatedContent(index int) Content {
	topics := []struct {
		category string
		topic    string
		badges   []string
	}{
		{"linux", "systemd service", []string{"linux", "systemd", "terminal"}},
		{"bash", "shell pipeline", []string{"bash", "cli", "automation"}},
		{"go", "Go worker", []string{"go", "testing", "concurrency"}},
		{"java", "Spring Boot endpoint", []string{"java", "spring-boot", "api"}},
		{"devops", "deployment runbook", []string{"devops", "ci", "reliability"}},
		{"tooling", "developer workstation", []string{"tooling", "linux", "productivity"}},
	}
	topic := topics[index%len(topics)]
	author := demoAuthors[index%len(demoAuthors)]
	created := time.Now().AddDate(0, 0, -(index%180 + 200)).Truncate(time.Second)
	published := created.Add(12 * time.Hour)
	return Content{
		Slug:        fmt.Sprintf("lab-note-%03d-%s", index+1, strings.ReplaceAll(topic.topic, " ", "-")),
		Title:       fmt.Sprintf("Lab Note %03d: %s", index+1, topic.topic),
		Summary:     fmt.Sprintf("A practical backlog note about a %s in the Urpi workshop.", topic.topic),
		Body:        fmt.Sprintf("This generated demo article explores a %s from a Linux-first engineering perspective.\n\nThe useful habit is to keep the command, the reason, and the rollback path close together. Measure the result, write down what changed, and leave the next reader a small example they can run.", topic.topic),
		Status:      "published",
		Category:    topic.category,
		Badges:      topic.badges,
		Author:      author,
		CreatedAt:   created,
		UpdatedAt:   created.Add(2 * time.Hour),
		PublishedAt: &published,
	}
}

// seedContent is the demo dataset, written in the voice of the backlog: notes on
// Linux, Go, Java, and the workflows around them.
var seedContent = []Content{
	{
		Slug:    "postmortems-that-people-read",
		Title:   "Postmortems That People Actually Read",
		Summary: "A blameless template short enough to be reused and specific enough to prevent a repeat.",
		Body: `A postmortem is a tool for changing behaviour, not a document for assigning blame. If nobody reads it twice, it did not work.

We keep the template to one page: what happened, who saw it first, what made it worse, what made it better, and the smallest change that would have prevented it. Longer analysis lives in linked notes so the summary stays skimmable.

Every postmortem ends with at most three follow-ups, each with an owner and a date. If a follow-up cannot be scheduled, we say so honestly instead of writing a wish list.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"reliability", "writing", "documentation"},
		CreatedAt: daysAgo(115), UpdatedAt: daysAgo(112),
	},
	{
		Slug:    "health-checks-and-readiness",
		Title:   "Health Checks and Readiness Are Different Questions",
		Summary: "Liveness asks whether the process survived; readiness asks whether it should receive traffic.",
		Body: `Conflating liveness and readiness turns a slow dependency into a restart loop.

A liveness probe answers one question: is this process stuck badly enough that restarting is the best option? It should not call the database, the queue, or anything else that can be slow without being fatal.

A readiness probe answers a different question: can this instance serve useful responses right now? That one is allowed to check dependencies, and it is allowed to fail temporarily. Failing readiness removes the instance from rotation; failing liveness restarts it. Different consequences deserve different probes.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"kubernetes", "reliability", "http"},
		CreatedAt: daysAgo(111), UpdatedAt: daysAgo(108),
	},
	{
		Slug:    "backups-you-have-restored",
		Title:   "Backups You Have Restored",
		Summary: "An untested backup is a rumour. A restore drill is evidence.",
		Body: `The only backup that counts is one you have restored, timed, and verified.

Once a quarter we take the most recent snapshot, restore it into a scratch environment, and check three things: how long it took, whether the row counts match expectations, and whether the application can start against it. The drill notes go next to the schedule, not in a separate wiki nobody opens.

If a restore needs a person who is on holiday, the backup is not finished. Fix the runbook before writing more backups.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"reliability", "postgresql", "automation"},
		CreatedAt: daysAgo(106), UpdatedAt: daysAgo(104),
	},
	{
		Slug:    "log-lines-worth-keeping",
		Title:   "Log Lines Worth Keeping",
		Summary: "Write logs for the person debugging at an inconvenient hour.",
		Body: `Logs are read under pressure, usually by someone who did not write the code and cannot change it quickly.

A useful line names the operation, the subject, and the outcome, and it includes the identifiers needed to find the related records. Redundancy across layers is noise: if the HTTP handler already logs the request, the repository should not log it again.

Levels are a contract. Debug is safe to delete, info describes normal progress, warn needs attention eventually, and error means a human should look. If everything is an error, nothing is.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"logging", "observability", "debugging"},
		CreatedAt: daysAgo(101), UpdatedAt: daysAgo(99),
	},
	{
		Slug:    "metrics-that-earn-their-keep",
		Title:   "Metrics That Earn Their Keep",
		Summary: "Four signals nobody has to interpret beat forty dashboards nobody opens.",
		Body: `Every metric costs something: collection, storage, and the attention needed to notice when it moves.

We start with latency, traffic, errors, and saturation for each user-facing path, then add a metric only when we can name the decision it would change. A dashboard that has never been opened during an incident is decoration.

Graphs are more useful with a line for the last release. Most arguments about a trend dissolve when the deployment markers are visible.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"metrics", "observability", "reliability"},
		CreatedAt: daysAgo(96), UpdatedAt: daysAgo(95),
	},
	{
		Slug:    "secrets-stay-out-of-git",
		Title:   "Secrets Stay Out of Git",
		Summary: "A short policy: no credentials in the repository, ever, including in history.",
		Body: `A secret committed once is a secret leaked, because history is copied, cached, and cloned.

The rule is simple: configuration lives in the repository, credentials do not. Local development reads from an untracked environment file, and pipelines read from the platform secret store. When a key must be rotated, we rotate it before deleting the commit that exposed it.

We also check for the accidental paste in review. A pre-commit hook that flags likely keys costs nothing and prevents the slowest kind of incident.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"encryption", "git", "tooling"},
		CreatedAt: daysAgo(91), UpdatedAt: daysAgo(90),
	},
	{
		Slug:    "blue-green-without-drama",
		Title:   "Blue-Green Without Drama",
		Summary: "Cutovers stay calm when both sides are deliberately boring.",
		Body: `Blue-green deployment is mostly an exercise in refusing to be clever.

Both environments run the same configuration, differ only in the image tag, and are exercised by the same smoke tests. The switch is a routing change that can be reverted in one step, which is the entire point: the rollback path is rehearsed because it is the same path used to deploy.

The difficult part is not the routing, it is the data. Any schema change must be compatible with both versions for the length of the cutover window, which is why we expand first, migrate, and contract later.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"reliability", "automation", "ci"},
		CreatedAt: daysAgo(86), UpdatedAt: daysAgo(85),
	},
	{
		Slug:    "cron-is-a-fine-scheduler",
		Title:   "Cron Is a Fine Scheduler",
		Summary: "Before reaching for a workflow engine, check whether cron plus a lock is enough.",
		Body: `Scheduled work is usually simpler than the tooling sold for it.

Cron is unglamorous and reliable. Add flock so overlapping runs cannot stack, write a log line with a run identifier, and make the job idempotent so a retry after a failure is harmless. That combination covers most nightly jobs we have ever needed.

Reach for a real orchestrator when the work has dependencies, partial failures that must be resumed, or a need for per-step visibility. Those are real requirements, but they should arrive before the tool does.`,
		Status:    "published",
		Category:  "devops",
		Badges:    []string{"cron", "automation", "linux"},
		CreatedAt: daysAgo(81), UpdatedAt: daysAgo(80),
	},
	{
		Slug:      "systemd-timers-are-just-cron-with-manners",
		Title:     "systemd Timers Are Just Cron With Manners",
		Summary:   "Timers bring logging, dependencies, and a useful status command to scheduled work.",
		Body:      `The first thing to say about systemd timers is that the schedule itself is not the interesting part. The interesting part is everything around it.\n\nA timer unit can declare that it runs after the network is up, that it jitters to avoid a thundering herd at midnight, and that failures are visible through the same status command used for services. Logs land in the journal with the rest of the machine.\n\nThe costs are also real: two unit files instead of one crontab line, and a syntax that rewards a quick look at the manual. For anything that matters, that trade is usually worth taking.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"systemd", "systemd-units", "automation"},
		CreatedAt: daysAgo(76), UpdatedAt: daysAgo(75),
	},
	{
		Slug:      "reading-a-cgroup-pressure-file",
		Title:     "Reading a cgroup Pressure File",
		Summary:   "PSI numbers explain whether a machine is slow because of CPU, memory, or IO.",
		Body:      `When a host feels slow, the usual averages hide the cause. Pressure stall information does not.\n\nThe files under the pressure directory report the share of time tasks were stalled waiting for a resource. A rising IO stall with idle CPU points somewhere very different from a saturated CPU, and the difference decides whether you look at disks or at code.\n\nThe numbers are also cumulative, so what matters is the rate of change. Read them twice, a minute apart, and the answer is usually obvious without installing anything.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"cgroups", "performance", "troubleshooting"},
		CreatedAt: daysAgo(71), UpdatedAt: daysAgo(70),
	},
	{
		Slug:      "lvm-for-people-who-like-undo",
		Title:     "LVM for People Who Like Undo",
		Summary:   "Logical volumes make storage changes reversible, which is the real feature.",
		Body:      `The headline feature of LVM is not snapshots, it is that resizing stops being frightening.\n\nA volume group turns fixed partitions into a pool, so growing a filesystem is a two-command operation and taking a snapshot before a risky change costs a few seconds. Snapshots are not backups, and they are not free, but they turn a mistake from an afternoon into a rollback.\n\nThe rule we follow is simple: any operation that could lose data gets a snapshot first, and the snapshot gets deleted the moment the change is verified.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"lvm", "filesystem"},
		CreatedAt: daysAgo(66), UpdatedAt: daysAgo(65),
	},
	{
		Slug:      "permissions-without-superstition",
		Title:     "Permissions Without Superstition",
		Summary:   "Read the octal number as three questions instead of memorising a table.",
		Body:      `Most permission confusion comes from treating the numbers as magic.\n\nRead them as three questions about the owner, the group, and everyone else: can they read, can they write, can they execute? Three bits each, four plus two plus one. Once that clicks, 755 and 644 stop being incantations.\n\nThe remaining decision is ownership. Files that a service must write belong to that service, not to a shared group with generous permissions. Narrow ownership and a clear group are easier to reason about than a permissive mode applied to everything.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"permissions", "filesystem"},
		CreatedAt: daysAgo(61), UpdatedAt: daysAgo(60),
	},
	{
		Slug:      "ssh-keys-you-can-audit",
		Title:     "SSH Keys You Can Audit",
		Summary:   "One key per machine, a comment that identifies it, and a rotation you have practised.",
		Body:      `A key you cannot identify is a key you cannot safely revoke.\n\nEvery key gets a comment naming the device and the year, lives in the agent with a passphrase, and is listed somewhere we can read during an incident. Shared keys are a shortcut that turns revocation into an archaeology project.\n\nThe rotation is deliberately dull: generate, install, confirm access, remove the old line. Doing it once while nothing is wrong is what makes it possible during an incident.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"ssh", "encryption", "tooling"},
		CreatedAt: daysAgo(56), UpdatedAt: daysAgo(55),
	},
	{
		Slug:      "wireguard-on-a-laptop",
		Title:     "WireGuard on a Laptop",
		Summary:   "A minimal tunnel config that reconnects quietly and stays out of the way.",
		Body:      `WireGuard is refreshingly small: a private key, a peer, an allowed range, and an endpoint.\n\nThe parts worth getting right are not the cryptography but the behaviour. Keep the allowed range tight so only the services you meant to reach are routed, enable the persistent keepalive only when you are behind a NAT, and let the tooling reconnect rather than adding a script that fights it.\n\nWe keep one configuration per purpose. A tunnel that tries to be the office, the lab, and the internet at once is the tunnel that breaks at the worst time.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"wireguard", "networking", "encryption"},
		CreatedAt: daysAgo(51), UpdatedAt: daysAgo(50),
	},
	{
		Slug:      "journalctl-in-anger",
		Title:     "journalctl in Anger",
		Summary:   "Filter by unit, by time, and by priority before you start scrolling.",
		Body:      `The journal is only intimidating until you use three flags.\n\nFilter to a unit, bound the time range to the incident window, and raise the priority floor so informational chatter disappears. Almost every investigation is these three filters plus a search for an identifier you already know.\n\nFollow mode is the other half of the tool: watching a unit while you reproduce a problem beats reconstructing the sequence afterwards from an unfiltered dump.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"logging", "systemd", "troubleshooting"},
		CreatedAt: daysAgo(46), UpdatedAt: daysAgo(45),
	},
	{
		Slug:      "keeping-updates-boring",
		Title:     "Keeping Updates Boring",
		Summary:   "Update often, reboot deliberately, and keep a note of what changed.",
		Body:      `The dangerous update is the one that has been postponed for months.\n\nWe apply updates on a schedule small enough that any single change is unremarkable, and we reboot at a moment we chose rather than one the machine chose. A short note in the log recording the date and the packages that moved is enough to connect a new symptom to a recent change.\n\nWhen something does break, the boring habit pays off: a small set of candidate changes is far easier to reason about than a year of accumulated drift.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"linux", "reliability", "troubleshooting"},
		CreatedAt: daysAgo(41), UpdatedAt: daysAgo(40),
	},
	{
		Slug:      "taming-grep-with-context",
		Title:     "Taming grep With Context",
		Summary:   "Three flags turn search output into something a human can read.",
		Body:      `Search results without context are a puzzle, not an answer.\n\nThe flags worth memorising are the ones that show neighbouring lines, print the filename, and count matches per file. Together they let you triage a repository without opening an editor. Add the fixed-string option and the recursive option and you rarely need anything else.\n\nThe habit that matters is searching for the thing that would have to be true for the bug to exist, rather than the message you already read in the logs.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"grep", "terminal", "debugging"},
		CreatedAt: daysAgo(36), UpdatedAt: daysAgo(35),
	},
	{
		Slug:      "a-tmux-layout-that-survives",
		Title:     "A tmux Layout That Survives",
		Summary:   "A layout per task, kept in a config file instead of in muscle memory.",
		Body:      `A terminal multiplexer earns its place when a session outlives a dropped connection.\n\nThe useful discipline is naming sessions after the task, not after the machine, and keeping a layout definition in the configuration so a rebuild looks the same as the original. Windows for editing, tests, and logs are enough for most work.\n\nThe extra benefit is psychological: an attached session makes it obvious what was in progress when you come back to a machine you last touched a week ago.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"tmux", "terminal", "tooling"},
		CreatedAt: daysAgo(131), UpdatedAt: daysAgo(130),
	},
	{
		Slug:      "makefiles-for-small-projects",
		Title:     "Makefiles for Small Projects",
		Summary:   "A dozen lines of make beat remembering six commands.",
		Body:      `Make is unfashionable and completely sufficient for the jobs a small repository actually has.\n\nWe keep targets for the things a newcomer will need in the first ten minutes: build, test, run, and clean. Declaring them phony avoids the classic confusion, and a single default target means the repository explains itself.\n\nThe test for whether a Makefile is doing too much is whether you can read it in one screen. When it grows past that, the build probably wants a real dependency graph, but most repositories never get there.`,
		Status:    "published",
		Category:  "linux",
		Badges:    []string{"makefile", "tooling", "automation"},
		CreatedAt: daysAgo(31), UpdatedAt: daysAgo(30),
	},
	{
		Slug:      "context-is-not-a-cancellation-token",
		Title:     "Context Is Not Just a Cancellation Token",
		Summary:   "Pass it as the first argument, cancel it at the boundary, and never store it.",
		Body:      `Context carries three things: cancellation, deadlines, and request-scoped values. Most bugs come from mixing them up.\n\nIt belongs in the function signature as the first parameter and nowhere else. Storing it in a struct makes lifetime ambiguous, and using it for optional configuration hides dependencies that should be arguments.\n\nThe cancellation should originate at the edge, where the request or the command decides how long the work may take. Every layer below just honours it. When a deadline expires, the useful question is which stage ignored it for too long.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"concurrency", "api", "errors"},
		CreatedAt: daysAgo(141), UpdatedAt: daysAgo(140),
	},
	{
		Slug:      "table-tests-that-still-read",
		Title:     "Table Tests That Still Read",
		Summary:   "A named case with a clear input, a clear expectation, and no shared state.",
		Body:      `A table test is only better than three separate tests if it stays readable.\n\nWe name every case with the behaviour it protects, keep setup inside the case where possible, and avoid clever closures that make failures hard to trace. When a case needs a paragraph of setup, it usually deserves to be its own function.\n\nThe failure message matters as much as the assertion. A test that prints the input, the actual value, and the expected value saves a debugging round trip in continuous integration.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"testing", "code-review"},
		CreatedAt: daysAgo(26), UpdatedAt: daysAgo(25),
	},
	{
		Slug:      "generics-when-they-help",
		Title:     "Generics, When They Help",
		Summary:   "Reach for type parameters when the algorithm is identical and only the type changes.",
		Body:      `Generics are a tool for removing duplicated algorithms, not for expressing clever abstractions.\n\nIf two functions differ only in the type they operate on, a type parameter removes the copy without hiding anything. If they differ in behaviour, an interface expresses the variation more honestly, because it names the operations that matter.\n\nThe cost is readability for readers who have not internalised the syntax, so we use them where the saving is obvious and skip them where an interface would document intent better.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"generics", "library", "architecture"},
		CreatedAt: daysAgo(21), UpdatedAt: daysAgo(20),
	},
	{
		Slug:      "error-wrapping-that-survives",
		Title:     "Error Wrapping That Survives",
		Summary:   "Add context once per layer, and let callers ask questions with errors.Is.",
		Body:      `An error message should read like a path from the entry point to the failure.\n\nWrap once per layer with the operation that failed, and keep the original error in the chain so callers can test for it. Repeating the same words at every level produces messages that are long and useless.\n\nSentinel errors answer a question about kind, custom types answer a question about detail, and everything else is context for a human. Deciding which of those a caller needs is the actual design work.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"errors", "debugging", "api"},
		CreatedAt: daysAgo(16), UpdatedAt: daysAgo(15),
	},
	{
		Slug:      "understanding-go-gc-pauses",
		Title:     "Understanding Go GC Pauses",
		Summary:   "Allocation rate, not heap size, is what usually drives pause behaviour.",
		Body:      `The collector is doing less work when the program allocates less.\n\nPauses grow when the allocation rate grows, because the collector has to run more often to keep up with the target. Reducing temporary allocations in a hot path often helps latency more than tuning any threshold.\n\nMeasure before tuning. The trace shows pause duration and frequency next to the work that caused it, which is the only way to tell an allocation problem apart from a scheduling one.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"gc", "memory", "performance"},
		CreatedAt: daysAgo(11), UpdatedAt: daysAgo(10),
	},
	{
		Slug:      "profiling-before-optimising",
		Title:     "Profiling Before Optimising",
		Summary:   "A profile turns an argument about performance into a list of the top three costs.",
		Body:      `Optimisation without a profile is guessing with extra steps.\n\nStart with the profile that matches the symptom: CPU for throughput, heap for memory, mutex and block profiles for contention. The first screen usually names a function nobody suspected.\n\nBenchmarks protect the change once it is made. A benchmark without a profile tells you that something got faster; a profile tells you why, which is what makes the next change obvious.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"profiling", "performance", "benchmarks"},
		CreatedAt: daysAgo(146), UpdatedAt: daysAgo(145),
	},
	{
		Slug:      "interfaces-small-enough-to-implement",
		Title:     "Interfaces Small Enough to Implement",
		Summary:   "Declare the interface where it is consumed, and keep it to the methods you call.",
		Body:      `A one-method interface is easy to satisfy and easy to read. A fifteen-method interface is a description of a package, not a seam.\n\nWe declare interfaces next to the code that consumes them, not next to the implementation. That keeps them honest: the only methods present are the ones a caller actually needs, which makes fakes trivial to write in tests.\n\nThe corollary is that adding a method to a widely used interface is a design change, not a convenience. It should come with a reason that mentions the caller.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"architecture", "api", "code-review"},
		CreatedAt: daysAgo(6), UpdatedAt: daysAgo(5),
	},
	{
		Slug:      "defer-in-a-loop",
		Title:     "defer in a Loop",
		Summary:   "Deferred calls run when the function returns, which is rarely what a loop wants.",
		Body:      `A deferred call inside a loop accumulates until the surrounding function ends.\n\nThat is usually harmless for a single mutex release and occasionally catastrophic for file handles or database rows. The fix is to move the body into its own function so each iteration has its own scope, which also tends to make the loop easier to read.\n\nWhen reviewing, missing deferral is one of the few patterns worth grepping for rather than reading for, because it is invisible in a quick scan.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"errors", "memory", "code-review"},
		CreatedAt: daysAgo(151), UpdatedAt: daysAgo(150),
	},
	{
		Slug:      "sync-pool-is-not-a-cache",
		Title:     "sync.Pool Is Not a Cache",
		Summary:   "It reduces allocation pressure; it does not promise that anything survives.",
		Body:      `A pool exists to shorten the life of short-lived objects, and it is allowed to drop everything at any moment.\n\nReaching for it to keep expensive things around is a mistake, because the runtime may clear the pool during a collection. If the resource must survive, it needs a real cache with a real lifetime policy.\n\nWhere the pool shines is in hot paths that allocate a buffer per request. The measured benefit is fewer collections, not faster individual operations.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"memory", "performance", "concurrency"},
		CreatedAt: daysAgo(156), UpdatedAt: daysAgo(155),
	},
	{
		Slug:      "json-encoding-details",
		Title:     "JSON Encoding Details That Bite",
		Summary:   "Field tags, zero values, and streaming are the three places bugs hide.",
		Body:      `Encoding looks like a solved problem until a zero value quietly disappears.\n\nThe omitempty option removes zero values, which is convenient for optional fields and surprising for booleans and counters that legitimately mean zero. Decide per field, and write the response shape down in a test so a future tag change is visible.\n\nFor large responses, stream into the writer instead of building the whole document in memory. It is barely more code and it removes a whole class of memory surprises.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"api", "http", "rest"},
		CreatedAt: daysAgo(161), UpdatedAt: daysAgo(160),
	},
	{
		Slug:      "graceful-shutdown",
		Title:     "Graceful Shutdown and the Requests You Forgot",
		Summary:   "Stop accepting new work, finish what is in flight, and give both a deadline.",
		Body:      `A process that exits immediately leaves a trail of half-finished work and confusing client errors.\n\nThe pattern is small: listen for the termination signal, stop accepting new connections, wait for in-flight requests with a deadline, and close resources in the reverse order they were opened. Background workers need the same treatment so their current item can finish.\n\nThe deadline is the important part. A shutdown with no bound is a shutdown that occasionally waits forever, which is worse than an abrupt stop.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"http", "reliability", "concurrency"},
		CreatedAt: daysAgo(166), UpdatedAt: daysAgo(165),
	},
	{
		Slug:      "reading-the-race-detector",
		Title:     "Reading the Race Detector",
		Summary:   "Two stacks, one shared address, and the write that proves the bug is real.",
		Body:      `A race report looks intimidating because it prints the same memory twice.\n\nThe first stack is one access, the second is the conflicting one, and the address they share names the variable. At least one of them is a write. Once you find that line, the question becomes which of the two goroutines owns the value, and the answer is usually that neither does.\n\nBecause the detector only reports races it actually observes, a clean run of the test suite is evidence, not proof. Keeping the flagged case as a regression test is what makes the fix durable.`,
		Status:    "published",
		Category:  "go",
		Badges:    []string{"concurrency", "testing", "debugging"},
		CreatedAt: daysAgo(171), UpdatedAt: daysAgo(170),
	},
	{
		Slug:      "virtual-threads-in-practice",
		Title:     "Virtual Threads in Practice",
		Summary:   "They remove the thread-per-request cost, not the need to bound work.",
		Body:      `Virtual threads make blocking code affordable again, which is a bigger change than it sounds.\n\nA request per virtual thread keeps the code sequential and readable while the runtime multiplexes onto a small pool of carrier threads. The catch is that a resource you block on still has to be bounded: a connection pool with twenty connections will queue no matter how cheap the threads are.\n\nThe awkward part is anything that pins the carrier, which keeps the blocking call stuck to a platform thread. Finding those cases is a matter of watching the thread dumps under load rather than guessing.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"threads", "concurrency", "performance"},
		CreatedAt: daysAgo(136), UpdatedAt: daysAgo(135),
	},
	{
		Slug:      "records-and-the-end-of-boilerplate",
		Title:     "Records and the End of Boilerplate",
		Summary:   "A record says what the data is, and the compiler writes the rest.",
		Body:      `Most value classes are a name and a list of fields, and records describe exactly that.\n\nThe generated accessors, equality, and string form remove a page of code that nobody reviewed anyway. Validation belongs in a compact constructor, which keeps the rule next to the fields it protects.\n\nThe remaining judgement is whether the type is genuinely a value. If it has behaviour, identity, or mutable state, a record is the wrong shape and a normal class documents the intent better.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"java", "api"},
		CreatedAt: daysAgo(226), UpdatedAt: daysAgo(225),
	},
	{
		Slug:      "gradle-builds-you-can-read",
		Title:     "Gradle Builds You Can Read",
		Summary:   "Keep the build script boring and put logic in tasks, not in the configuration block.",
		Body:      `A build that needs a debugger is a build nobody will maintain.\n\nWe keep plugins and versions in one catalog, declare dependencies in one obvious place, and avoid clever logic in configuration. Anything conditional belongs in a task with a name, so it appears in the task list and can be run on its own.\n\nBuild time is a feature. Configuration cache and a warm daemon are usually worth more than any micro-optimisation inside a task.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"gradle", "tooling", "ci"},
		CreatedAt: daysAgo(221), UpdatedAt: daysAgo(220),
	},
	{
		Slug:      "jvm-memory-you-can-explain",
		Title:     "JVM Memory You Can Explain",
		Summary:   "Heap, metaspace, and the stack each answer a different question.",
		Body:      `Memory problems are easier to reason about once the three areas are distinct.\n\nThe heap holds objects and is what most tuning is about. Metaspace holds class metadata, which is why a classloader leak looks like native memory growth. The stack is per thread, which is why a thousand threads cost more than the objects they touch.\n\nReading a heap histogram before and after a suspected leak answers the question faster than any amount of configuration tuning, because it names the class that is accumulating.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"jvm", "memory", "profiling"},
		CreatedAt: daysAgo(216), UpdatedAt: daysAgo(215),
	},
	{
		Slug:      "checked-exceptions-revisited",
		Title:     "Checked Exceptions Revisited",
		Summary:   "They are a documentation tool, and they are easy to abuse.",
		Body:      `A checked exception forces a caller to acknowledge that something can fail, which is genuinely useful for recovery paths.\n\nThe abuse is a signature that lists five checked exceptions, most of which nobody can act on. Those belong in a single application exception or in an unchecked type, so the interesting failure stays visible.\n\nThe honest test is whether a caller would do something different for this failure. If not, it should not be in the signature.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"errors", "java", "api"},
		CreatedAt: daysAgo(211), UpdatedAt: daysAgo(210),
	},
	{
		Slug:      "streams-are-not-always-faster",
		Title:     "Streams Are Not Always Faster",
		Summary:   "They express a pipeline; they do not promise a faster one.",
		Body:      `A stream describes what should happen, and the runtime decides how.\n\nThat trade is usually worth it for readability, especially for filter-and-map work over a collection. It is a poor trade in the hottest loop of a program, where an explicit loop with a preallocated array is often measurably faster and no less clear.\n\nThe functional style also hides cost when intermediate operations allocate. If the pipeline is in a hot path, measure it with the same discipline used for anything else.`,
		Status:    "published",
		Category:  "java",
		Badges:    []string{"performance", "java", "benchmarks"},
		CreatedAt: daysAgo(206), UpdatedAt: daysAgo(205),
	},
	{
		Slug:      "set-e-is-not-a-strategy",
		Title:     "set -e Is Not a Strategy",
		Summary:   "Errexit helps, but it does not replace checking what you actually care about.",
		Body:      `Turning on errexit catches the common case and gives a false sense of safety.\n\nCommands used as conditions, anything in a pipeline except the last element, and failures inside a subshell can all slip past. Reaching for the pipefail option fixes one of those and leaves the rest.\n\nThe scripts that stay reliable are the ones that check the result of anything destructive, print what they are about to do, and are short enough to read in one sitting. Options help; they do not decide.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"bash", "errors", "automation"},
		CreatedAt: daysAgo(196), UpdatedAt: daysAgo(195),
	},
	{
		Slug:      "quoting-is-the-whole-language",
		Title:     "Quoting Is the Whole Language",
		Summary:   "Almost every mysterious shell bug is a missing pair of quotes.",
		Body:      `The shell splits, expands, and globs, and quoting is how you say no.\n\nDouble quotes keep a value together while still expanding variables. Single quotes keep everything literal. Unquoted values are where filenames with spaces become two arguments and an empty variable disappears entirely.\n\nThe habit worth building is quoting every expansion by default and removing quotes only when word splitting is genuinely wanted, which is rare and should look deliberate.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"bash", "shell", "debugging"},
		CreatedAt: daysAgo(191), UpdatedAt: daysAgo(190),
	},
	{
		Slug:      "pipes-without-surprises",
		Title:     "Pipes Without Surprises",
		Summary:   "Buffering, exit codes, and missing input are the three things to watch.",
		Body:      `A pipeline is a small distributed system, and it fails in the same interesting ways.\n\nThe exit status of a pipeline is the status of its last command unless you ask for more. A consumer that exits early leaves the producer writing into a closed pipe. Buffering means output can appear in a different order than you expect when watching a live log.\n\nWriting the intermediate result to a file while developing is not a defeat. It turns three commands into two understandable steps and usually finds the bug faster.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"bash", "shell", "troubleshooting"},
		CreatedAt: daysAgo(186), UpdatedAt: daysAgo(185),
	},
	{
		Slug:      "arrays-in-bash",
		Title:     "Arrays in Bash, Carefully",
		Summary:   "Quote the expansion, iterate over values, and do not use strings as lists.",
		Body:      `Arrays exist so that filenames with spaces stop being a problem.\n\nExpanding an array unquoted reintroduces exactly that problem, so the expansion gets quotes like any other value. Iterating over the elements rather than over an index keeps the loop simple and avoids off-by-one mistakes.\n\nThe alternative, packing a list into a delimited string, works until the first value contains the delimiter. If the list can contain arbitrary text, use the array.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"bash", "shell", "snippets"},
		CreatedAt: daysAgo(181), UpdatedAt: daysAgo(180),
	},
	{
		Slug:      "trap-and-cleanup",
		Title:     "trap and Cleanup",
		Summary:   "Temporary files and locks should be released even when the script fails.",
		Body:      `A script that leaves state behind is a script that breaks the next run.\n\nA trap on exit runs whether the script succeeded, failed, or was interrupted, which makes it the right place to remove temporary directories and release locks. Handling the interrupt signal lets a long job stop between steps instead of in the middle of one.\n\nThe cleanup should be idempotent, because a trap can run when the resource was never created. Checking before removing costs one line and prevents a confusing error at exactly the worst moment.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"bash", "automation", "reliability"},
		CreatedAt: daysAgo(176), UpdatedAt: daysAgo(175),
	},
	{
		Slug:      "find-with-print0",
		Title:     "find With print0",
		Summary:   "Newline-separated filenames break on real filesystems.",
		Body:      `Passing a list of paths through a newline is an assumption that filenames are well behaved.\n\nUsing the null-terminated output with the matching input option removes the assumption entirely, which matters as soon as a filename contains a space, a quote, or a newline. It costs one extra flag on each side.\n\nThe other half of the habit is limiting the search early: by depth, by type, or by modification time. A find that scans an entire home directory to produce three results is a find that should have been narrower.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"filesystem", "bash", "terminal"},
		CreatedAt: daysAgo(121), UpdatedAt: daysAgo(120),
	},
	{
		Slug:      "awk-one-liners-worth-knowing",
		Title:     "awk One-Liners Worth Knowing",
		Summary:   "Choosing fields, filtering rows, and summing a column cover most real work.",
		Body:      `Awk is a tiny language for columns, and three of its features handle most tasks.\n\nPrinting selected fields, filtering lines by a condition, and accumulating a total are enough to replace a great deal of fragile text processing. Field numbering is one-based, and treating the first line differently is a one-line idiom.\n\nThe reason to prefer awk over a chain of substitutions is that it parses structure rather than matching text, which makes it survive changes that would break a regular expression.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"awk", "terminal", "snippets"},
		CreatedAt: daysAgo(116), UpdatedAt: daysAgo(114),
	},
	{
		Slug:      "sed-substitutions-that-are-safe",
		Title:     "sed Substitutions That Are Safe",
		Summary:   "Preview, choose an unusual delimiter, and back up before writing in place.",
		Body:      `An in-place substitution without a backup is a one-way door.\n\nThe safe habit is to run the substitution first and read the output, then repeat it with the in-place flag and a backup suffix. Choosing a delimiter that does not appear in the pattern removes a whole category of quoting confusion when the text involves paths.\n\nFor anything beyond a simple substitution, the honest answer is usually that the task wants a real parser. Sed is excellent at what it does and quietly dangerous outside that.`,
		Status:    "published",
		Category:  "bash",
		Badges:    []string{"sed", "regex", "terminal"},
		CreatedAt: daysAgo(111), UpdatedAt: daysAgo(109),
	},
	{
		Slug:      "docker-images-that-stay-small",
		Title:     "Docker Images That Stay Small",
		Summary:   "Multi-stage builds, static binaries, and a base image with nothing extra.",
		Body:      `Most image size comes from build tooling that should not exist at runtime.\n\nA multi-stage build keeps the compiler, the package manager, and the source tree in the first stage and copies only the finished artifact into the second. For a language that can produce a static binary, the runtime stage needs nothing but certificates and a non-root user.\n\nThe payoff is not only size. A smaller image starts faster, has fewer packages to patch, and makes the container's purpose obvious to anyone reading the file.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"docker", "containers", "ci"},
		CreatedAt: daysAgo(106), UpdatedAt: daysAgo(105),
	},
	{
		Slug:      "dockerfiles-people-can-trust",
		Title:     "Dockerfiles People Can Trust",
		Summary:   "Pin versions, order layers by change frequency, and run as a normal user.",
		Body:      `A Dockerfile is a build script with a cache, and both parts deserve attention.\n\nPinning the base version makes a rebuild reproducible instead of a surprise. Copying the dependency manifest and installing before copying the source keeps the expensive layer cached while the code changes underneath it. Creating an unprivileged user costs two lines and removes a whole class of escalation.\n\nWe also like a build that fails loudly: no silent ignoring of errors, no downloads without checksums, no assumptions about the network.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"docker", "containers", "encryption"},
		CreatedAt: daysAgo(206), UpdatedAt: daysAgo(204),
	},
	{
		Slug:      "compose-for-local-only",
		Title:     "Compose for Local Only",
		Summary:   "A good compose file documents how to run the project on a laptop.",
		Body:      `The value of a compose file is onboarding, not orchestration.\n\nIt should start the one or two dependencies a developer actually needs, with health checks so the app does not race them, and named volumes so data survives a restart. Anything that only matters in production belongs somewhere else.\n\nKeeping the file short is a feature. If a newcomer has to read it for five minutes, it has stopped doing its job.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"docker", "compose", "tooling"},
		CreatedAt: daysAgo(101), UpdatedAt: daysAgo(100),
	},
	{
		Slug:      "changelogs-with-a-point-of-view",
		Title:     "Changelogs With a Point of View",
		Summary:   "Write for the reader deciding whether to upgrade.",
		Body:      `A changelog is a decision aid, not a commit log.\n\nEntries grouped by what changed for the user, with breaking changes called out first and the reason stated briefly, are worth more than a complete list of tickets. Linking to the relevant documentation turns an announcement into something actionable.\n\nThe tone matters too. A changelog that says only what was fixed leaves readers wondering what is still broken. Naming the trade-offs builds the trust that makes the next release easier to adopt.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"documentation", "writing", "changelog"},
		CreatedAt: daysAgo(96), UpdatedAt: daysAgo(95),
	},
	{
		Slug:      "linters-earn-their-place",
		Title:     "Linters Have to Earn Their Place",
		Summary:   "Every rule should prevent a bug somebody actually shipped.",
		Body:      `A linter is a code review that never gets tired, and like a code review it can be wrong.\n\nWe enable rules with a reason attached: this one catches a mistake we made twice, that one keeps a style argument from recurring. Rules that only express a preference are turned off deliberately rather than tolerated forever.\n\nFormatting belongs to a formatter, not to a discussion. Agreeing to stop arguing about whitespace is one of the cheapest productivity gains available.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"tooling", "code-review", "ci"},
		CreatedAt: daysAgo(91), UpdatedAt: daysAgo(90),
	},
	{
		Slug:      "git-bisect-when-things-break",
		Title:     "git bisect When Things Break",
		Summary:   "A scripted good and bad commit turns a mystery into a short search.",
		Body:      `Bisect is the fastest way to answer "when did this start?" and it is usually forgotten.\n\nMarking a known good and a known bad revision halves the search space each step, which finds the culprit in a handful of builds even in a long history. Automating the check with a script makes it possible to walk away while it runs.\n\nThe prerequisite is a test that fails reliably. If the check is flaky, bisect returns a confident wrong answer, which is worse than no answer at all.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"git", "debugging", "testing"},
		CreatedAt: daysAgo(86), UpdatedAt: daysAgo(85),
	},
	{
		Slug:      "editor-config-as-a-document",
		Title:     "Editor Config as a Document",
		Summary:   "Aligning the small conventions removes a recurring review comment.",
		Body:      `Line endings, indentation, and trailing whitespace produce a surprising share of review noise.\n\nDeclaring them once in a shared file means the editor agrees with continuous integration and with colleagues on other platforms. It is a small file with a large return, because it removes a class of comment that carries no information.\n\nThe same idea applies to the manifest that pins the toolchain version. When the version is written down, works-on-my-machine becomes a conversation about the file rather than about machines.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"tooling", "code-review", "documentation"},
		CreatedAt: daysAgo(79), UpdatedAt: daysAgo(78),
	},
	{
		Slug:      "vim-motions-worth-practising",
		Title:     "Vim Motions Worth Practising",
		Summary:   "Structure-aware movement beats repeating a character search.",
		Body:      `The motions that pay back their practice time are the ones that move by structure.\n\nJumping to a word, to a paragraph, to a matching bracket, and to the last edit covers most navigation. Combining a motion with an operator is what turns editing into a sentence rather than a sequence of keystrokes.\n\nReaching for the search command instead of holding a movement key is the single habit that changes the most, because it turns a long scroll into one jump.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"vim", "terminal", "tooling"},
		CreatedAt: daysAgo(76), UpdatedAt: daysAgo(75),
	},
	{
		Slug:      "http-clients-that-timeout",
		Title:     "HTTP Clients That Time Out",
		Summary:   "The default has no deadline, which is a bug waiting for a slow dependency.",
		Body:      `A client without a timeout will wait as long as the server feels like making it wait.\n\nSetting a total timeout, a connection timeout, and a response header timeout separately makes the failure mode easy to reason about. Retries belong on top of those, bounded and only for idempotent requests, ideally with a little randomised delay.\n\nThe same discipline applies to the server side. A handler that can block forever turns one slow dependency into an outage, and a deadline propagated from the request is what prevents it.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"http", "reliability", "api"},
		CreatedAt: daysAgo(176), UpdatedAt: daysAgo(174),
	},
	{
		Slug:      "yaml-is-a-data-format-honestly",
		Title:     "YAML Is a Data Format, Honestly",
		Summary:   "Two spaces, no tabs, and quote anything that looks like a special value.",
		Body:      `Most configuration mystery comes from YAML being more clever than expected.\n\nIndentation is significant and tabs are forbidden, so a consistent editor setting prevents the most common failure. Values that look like booleans or numbers are coerced unless quoted, which is how a version number becomes a float.\n\nThe practical rule is to treat YAML as JSON with comments: keep the structure shallow, quote the ambiguous scalars, and let a validator catch the rest before a deployment does.`,
		Status:    "published",
		Category:  "tooling",
		Badges:    []string{"yaml", "tooling", "troubleshooting"},
		CreatedAt: daysAgo(171), UpdatedAt: daysAgo(170),
	},
	{
		Slug:      "writing-for-readers-in-a-hurry",
		Title:     "Writing for Readers in a Hurry",
		Summary:   "The first sentence should answer the question that brought them here.",
		Body:      `Most readers arrive with a question and very little patience.\n\nPutting the answer in the first sentence, before the context and before the caveats, respects that. Background belongs after the answer, where a reader who needs it will happily go.\n\nThe structure that survives editing is: what, so what, now what. What happened, why it matters to you, and what you should do next. Everything else is optional.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"writing", "editorial", "documentation"},
		CreatedAt: daysAgo(166), UpdatedAt: daysAgo(165),
	},
	{
		Slug:      "headings-are-a-table-of-contents",
		Title:     "Headings Are a Table of Contents",
		Summary:   "Read only the headings; if the outline does not make sense, the piece does not either.",
		Body:      `A reader skimming the headings should be able to reconstruct the argument.\n\nThat test catches most structural problems before a word of prose is rewritten. If two headings say the same thing, there is one section too many. If a heading is a noun phrase with no verb, it probably describes a topic instead of a claim.\n\nWe keep one idea per heading level and avoid skipping levels, because the hierarchy is a promise about how the parts relate.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"writing", "editorial", "typography"},
		CreatedAt: daysAgo(161), UpdatedAt: daysAgo(160),
	},
	{
		Slug:      "the-first-draft-is-for-you",
		Title:     "The First Draft Is for You",
		Summary:   "Write it badly, then delete the parts that only helped you think.",
		Body:      `The first draft exists to discover what the piece is actually about.\n\nThat means it is allowed to be repetitive, uncertain, and too long. The editing pass is where the discovered argument gets stated once and the scaffolding comes down.\n\nThe tell-tale sign of an unedited draft is explanation that answers a question the reader never asked. It was necessary to write; it is not necessary to read.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"writing", "editing", "editorial"},
		CreatedAt: daysAgo(156), UpdatedAt: daysAgo(155),
	},
	{
		Slug:      "accessibility-is-not-a-checklist",
		Title:     "Accessibility Is Not a Checklist",
		Summary:   "Contrast and keyboard access are the floor, not the ceiling.",
		Body:      `Passing an automated audit is the beginning of accessibility work, not the end.\n\nThe parts a tool cannot see are usually the parts that matter most: whether the reading order makes sense, whether a change is announced to assistive technology, whether an error explains how to recover. Those come from using the thing the way a person with different needs would.\n\nThe cheapest improvement in any interface is almost always better contrast and a visible focus state. Both are one line of styling and both help everybody.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"accessibility", "css", "design"},
		CreatedAt: daysAgo(151), UpdatedAt: daysAgo(150),
	},
	{
		Slug:      "empty-states-are-real-states",
		Title:     "Empty States Are Real States",
		Summary:   "The screen with no data is the first screen many people ever see.",
		Body:      `An empty list is not an edge case; for a new account it is the default.\n\nA good empty state says why there is nothing here, what would put something here, and offers the action that does it. It is the most useful onboarding surface in the product and it is usually left as a blank rectangle.\n\nThe same applies to loading and error states. If those three were designed alongside the populated view instead of after it, most interface awkwardness would disappear.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"design", "css", "typography"},
		CreatedAt: daysAgo(146), UpdatedAt: daysAgo(145),
	},
	{
		Slug:      "boring-typography-scales",
		Title:     "Boring Typography Scales",
		Summary:   "Two families, a handful of sizes, and a generous line height.",
		Body:      `Typography decisions get easier when the scale is small enough to remember.\n\nOne family for prose and one for code, four or five sizes total, and a line height that grows as the text gets smaller. Constraining the choices removes the meeting about which heading deserves an extra two pixels.\n\nThe most common mistake is not a bad font, it is too little spacing. Reading is easier with slightly more space than looks correct on a designer's monitor.`,
		Status:    "published",
		Category:  "design",
		Badges:    []string{"typography", "css", "design"},
		CreatedAt: daysAgo(141), UpdatedAt: daysAgo(140),
	},
	{
		Slug:      "deep-work-and-the-backlog",
		Title:     "Deep Work and the Backlog",
		Summary:   "A backlog that everyone can read is a backlog that interrupts everyone.",
		Body:      `A shared list of everything is a shared invitation to interrupt.\n\nWe keep a small, visible queue of committed work and a longer list that nobody is expected to watch. The visible part is what the team is doing now, and it should be short enough to hold in mind.\n\nProtecting a block of uninterrupted time is not a personal preference, it is a scheduling decision with a cost attached. Naming the cost makes it possible to defend.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"roadmap", "workflow", "writing"},
		CreatedAt: daysAgo(136), UpdatedAt: daysAgo(135),
	},
	{
		Slug:      "code-review-as-a-conversation",
		Title:     "Code Review as a Conversation",
		Summary:   "Ask about intent before style, and separate blocking from optional.",
		Body:      `A review that begins with formatting has skipped the part where the author learns something.\n\nThe useful order is intent, then correctness, then clarity, then style. The first question is whether the change does what it claims; everything else is cheaper to discuss once that is settled.\n\nLabelling comments as blocking or optional removes the guesswork. A reviewer who writes only the blocking ones is being kind, not lazy.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"code-review", "review", "writing"},
		CreatedAt: daysAgo(131), UpdatedAt: daysAgo(130),
	},
	{
		Slug:      "documentation-next-to-the-code",
		Title:     "Documentation Next to the Code",
		Summary:   "A document in the repository is reviewed, versioned, and found.",
		Body:      `Documentation that lives only in a wiki drifts out of date without anyone noticing.\n\nKeeping the explanation in the repository means a change that invalidates it appears in the same review as the change itself. Comments explain why, the readme explains how to run it, and the design notes explain what was rejected.\n\nThe best documentation is short, current, and located where somebody will trip over it. A perfect document nobody opens is a hobby.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"documentation", "writing", "code-review"},
		CreatedAt: daysAgo(126), UpdatedAt: daysAgo(125),
	},
	{
		Slug:      "estimates-as-ranges",
		Title:     "Estimates as Ranges",
		Summary:   "A range with an assumption attached is more useful than a date without one.",
		Body:      `A single date invites an argument about the date rather than the work.\n\nA range plus the assumption it rests on gives the reader something to reason with: if the assumption holds, the work lands here, and if the integration turns out to be more involved, it moves. That is a more honest sentence than any point estimate.\n\nThe practice also surfaces hidden uncertainty early. When the range is very wide, the useful next step is usually a small experiment rather than more discussion.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"roadmap", "workflow", "writing"},
		CreatedAt: daysAgo(121), UpdatedAt: daysAgo(120),
	},
	{
		Slug:      "turning-incidents-into-habits",
		Title:     "Turning Incidents Into Habits",
		Summary:   "A fix in a document gets forgotten; a fix in automation does not.",
		Body:      `The end of an incident is the moment when a lesson is easy to capture and easy to lose.\n\nThe durable outcome is not the write-up, it is the check that runs automatically, the alert that fires earlier, or the guard rail that makes the wrong action impossible. Prefer the change that works when nobody remembers the incident.\n\nWe keep one follow-up of that kind per incident. A list of ten improvements usually produces zero, and one automation produces a permanent reminder.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"reliability", "automation", "reminders"},
		CreatedAt: daysAgo(116), UpdatedAt: daysAgo(115),
	},
	{
		Slug:      "shipping-the-boring-improvement",
		Title:     "Shipping the Boring Improvement",
		Summary:   "The change nobody will notice is often the one that removes the most pain.",
		Body:      `Rewrites are exciting and rarely the best use of a week.\n\nThe improvements that hold up are dull: a clearer error message, one fewer manual step, a test that closes a flaky path, a log line that names the request. None of them make a good demo, and together they change how the work feels.\n\nWe protect time for them explicitly. If only the interesting work gets scheduled, the boring improvements accumulate as complaints instead of commits.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"workflow", "review", "roadmap"},
		CreatedAt: daysAgo(111), UpdatedAt: daysAgo(110),
	},
	{
		Slug:      "notes-from-automating-my-own-workflow",
		Title:     "Notes From Automating My Own Workflow",
		Summary:   "Automate the second time you do something, not the first.",
		Body:      `Automating a task you have done once produces a script shaped like a guess.\n\nThe second time, the variation is visible: which parts are always the same and which parts change. That is the moment a script earns its keep, because the interface it needs is now obvious.\n\nThe other rule is to keep the automation readable. A script with no output and no comments saves five minutes and costs thirty the next time it fails silently.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"automation", "tooling", "bash"},
		CreatedAt: daysAgo(106), UpdatedAt: daysAgo(105),
	},
	{
		Slug:      "the-meeting-that-should-be-a-note",
		Title:     "The Meeting That Should Be a Note",
		Summary:   "If nobody will ask a question, write it down instead.",
		Body:      `Recurring meetings survive long after the reason for them has been forgotten.\n\nThe test is simple: is there a decision to make or a question to answer that requires people in the same moment? If the answer is that the meeting exists to share status, a short written update does the job better and consumes less of everyone's attention.\n\nConverting one meeting into a note is a cheap experiment. If nothing breaks, the meeting was already finished.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"workflow", "writing", "documentation"},
		CreatedAt: daysAgo(66), UpdatedAt: daysAgo(65),
	},
	{
		Slug:    "built-with-ai",
		Title:   "This Demo Was Built With AI, and That Is the Point",
		Summary: "This whole backlog, the API behind it, and the image it runs in were written in conversation with an AI assistant. Here is what that felt like.",
		Body: `This application is a demo, and it has a confession: almost all of it was written together with an AI assistant.

The backend is a small Go service with a hand-written router, an in-memory store, and one dependency-free security seam. The frontend is React and TypeScript. The whole thing ships as a single container with the compiled binary and the static assets inside, and the articles you are reading are seeded from a Go slice rather than a database. None of that was typed out from a blank page. It was described, generated, reviewed, corrected, and described again.

What worked well was the speed of the boring parts. Renaming a field across a layer, adding a test that mirrors an existing one, restructuring a file so the pieces sit in a sensible order, writing fifty articles in a consistent voice: these are the tasks where a clear instruction produces a usable result in seconds. The assistant was also good at the thing people often forget to ask for, which is consistency. Once it had seen the shape of the codebase, it kept producing code that looked like it belonged there.

What worked less well was anything that depended on judgement that had not been stated out loud. Left alone, the assistant will happily keep a database that nobody asked for, add an abstraction layer to a five-line function, or write a test that asserts nothing meaningful. Every one of those happened during this build. The fix was never clever prompting, it was reading the diff.

That is the habit this project mostly taught. AI is excellent at producing plausible code quickly, and plausible is not the same as correct. The security model in this repository changed shape twice: once when it became clear that a demo with one owner should not carry an account system at all, and once when the storage layer was replaced rather than repaired. Both decisions were human, and both were simpler than what had grown there before. The assistant implemented them faithfully and quickly; it did not propose them.

A few practical notes, in case this is useful to someone reading it as a case study rather than as an article.

Describe the constraint, not just the task. Asking for a single self-contained image produced a very different result from asking for a deployment.

Let the assistant see the whole picture. The quality of generated code tracked the accuracy of the surrounding context almost exactly.

Insist on verification. Every claim in this project was checked by building it, testing it, and running it. The assistant is at least as good at explaining why something works as it is at making it work, so the explanation was never treated as evidence.

Keep the changelog honest. The record of what was asked and what was done is more interesting than the result, and it is the only part of this that a future reader cannot reconstruct.

So the demo is small, and it is finished, and an AI wrote most of the words and most of the lines. The architecture, the trade-offs, and the decision to throw away the database were mine. That division of labour is the interesting part, and it is probably how a lot of small software gets built now. It was, honestly, fun.`,
		Status:    "published",
		Category:  "workflow",
		Badges:    []string{"open-source", "writing", "documentation"},
		CreatedAt: daysAgo(2), UpdatedAt: daysAgo(1),
	},
}
