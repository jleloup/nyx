---
title: Build and Automation
layout: single
toc: true
permalink: /guide/user/best-practice/build-and-automation/
---

Nowadays all teams use CI/CD platforms to automate the build and release process around their software projects but often times they get into some pitfalls like:

* treating the build process as a second class citizen, falling short in proper engineering and maintenance
* introducing discrepancies between the CI/CD scripts and local scripts used by developers
* formalizing only parts of the process

These topics fall under the *DevOps* discipline (or *DevSecOps*) and that involves a broad range of subjects outside the scope of this page. On the other hand we'll just focus on a few specific points here as follows.

## Build scripts

The entire build process should be **consistent**, **portable** and **efficient**.

A **consistent** process doesn't need to repeat the same tasks in order to grant that the output of one task never changes provided the same input. Single tasks have their dependencies modelled in order to balance granularity and reuse of already available artifacts but under no circumstance the same task should be repeated to produce the same artifacts. If you find yourself in duplicating code in the build scripts you likely need to refactor task dependencies or parametrize the script.

To be **portable**, the build proces can be executed on all the platforms and environments used by the team yielding to the same outcomes **idempotently**.

An **efficient** build process avoids consuming unnecessary resources and takes advantage of **incremental** builds and parallel execution to shorten the time needed to complete.

How does all this relate to Nyx and the release process?

In the first place the release process (and the tools it uses) must not break any of the above principles. In particular the release management tool must not break the **consistency** principle (as many of the tools out there do). The *release* task should be **incremental** but most release management tools encapsulate all of their logic in one single *atomic* task and, for example, do not allow you to know the version number before the *release* task, which usually runs in the end and only in some circumstances. The issue here is that:

* all of the previous tasks (like *build*, *test*, *publish* etc) need to generate their own *fictional* version number to build and package artifacts but that's not the same number that the *release* task will use in the end so, when releasing, the generation of artifacts must be executed again to use the right version. To narrow the negative impact of repeating tasks you do not re-run tests or, if you do, you lose all of the **efficiency** as you basically end up by repeating the entire process by dependencies of the *release* task
* in those configurations where the *release* task must not run (i.e. in *feature branches*) you can't use the version number generated by the build tool so your build process is not **portable** and really reliable

Nyx solves the above issues as it's **incremental** as its internal workflow is split in two phases: the first one ([**inference**]({{ site.baseurl }}{% link _pages/guide/user/02.introduction/how-nyx-works.md %}#infer)) only fetches and elaborates informations (like the past and next version numbers), while the second one ([**publication**]({{ site.baseurl }}{% link _pages/guide/user/02.introduction/how-nyx-works.md %}#publish)), usually executed last in the build process, actually issues the release. In case the *release* task doesn't run (i.e. because the process breaks because of failed tests or the current configuration or branch doesn't require to release), Nyx just ensures that the version that was generated is still **consistent** so an **incremental** build can start from where it left.

### Task Order

In a standard build script you probably have the *publish* and the *release* tasks, among others. You shouldn't collapse them into one just because they have different objectives.

The *publish* task uploads artifacts and makes them available to your audience and this may happen also for *internal* releases. The *release* task, on the other hand, performs some final activities to *seal* the release and make it *official*. The *release* task should not be used to *publish* anything with the only exception of the release manifest itself (often the *changelog*).

In summary, you should keep the *publish* and *release* steps separate and run *release* after all the previous tasks have completed successfully.

Whether to run the *publish* before or after the *release* step is often controversial. We believe that *release* should be executed as the very final step because while you can publish additional releases (in case you can't *roll back* a publication by deleting the artifacts), releasing multiple times might be very misleading for the project consumers. From a semantic standpoint the *release* is also a sort of *stamp* that, once applied, grants that the version being delivered has passed all of the required steps, including publications.
{: .notice--info}

## Environment Segregation

The *release* task should be considered critical as it affects the audience and consumers of your project. Issuing wrong releases or not granting a consistent release process has serious impact on the reputation of your project and that's why you should release carefully.

Release should only be issued from a centralized CI/CD environment with limited access for configuration changes and only after all the tests have been performed. Releases should never be issued from local developer environments and should never overwrite an old release once it has been published.

In order to enforce this principle just make sure that the credentials needed by the release process are kept secret and only available on the CI/CD environment, usually as environment variables. Never share those credentials with anyone nor store them in any file within the repository.

Nyx gives you the tools make steps of the release process conditional, depending on the environment it's running on, so you can achieve full segregation dynamically.

## CI/CD vs local scripts

Chances are that when using a CI/CD platform you also need to build locally before you commit so your scripts need to be **portable** to run **idempotently** on both environments to ensure that the process itself is defined and tested just like any other artifact in your project. It also means that the script doesn't make strong assumptions on one environment or the other.

To achieve all this you should author your primary scripts using a build tool of your choice (i.e. Gradle, Maven etc) and have the CI/CD platform invoke those scripts, instead duplicating the logic. In other words you treat your primary build scripts as the *single source of truth* for your build process while CI/CD platforms may have additional scripts to bridge the gap between the hosting platform and the primary script, when necessary.