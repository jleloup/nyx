---
title: State Reference
layout: single
toc: false
permalink: /guide/user/state-reference/
---

You can get access to the internal state of Nyx if you need to use one of the values it computes internally or just to inspect its internals, like for debugging purpose.

You can save the State by setting the [`stateFile`]({{ site.baseurl }}{% link _pages/guide/user/03.configuration-reference/global-options.md %}#state-file) configuration option and even use it to break the execution in phases, in which case you can load a previously stored State file by enabling the [`resume`]({{ site.baseurl }}{% link _pages/guide/user/03.configuration-reference/global-options.md %}#resume) option. This way, when working on large repositories, where the release process may take long to complete, you can split the process in phases (and run other tasks in between), making sure that you always start from where you left.

You can access the State attributes directly from **Gradle** build scripts as explained [here]({{ site.baseurl }}{% link _pages/guide/user/02.introduction/usage.md %}#accessing-the-nyx-state-extra-project-property-from-build-scripts).
{: .notice--info}

In this section you can find details about every value contained in the State file.