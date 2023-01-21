/*
 * Copyright 2020 Mooltiverse
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
This is the Git package for Nyx, encapsulating the underlying Git implementation.
*/
package git

import (
	"fmt"     // https://pkg.go.dev/fmt
	"strings" // https://pkg.go.dev/strings

	ggit "github.com/go-git/go-git/v5"                             // https://pkg.go.dev/github.com/go-git/go-git/v5
	ggitconfig "github.com/go-git/go-git/v5/config"                // https://pkg.go.dev/github.com/go-git/go-git/v5
	ggitplumbing "github.com/go-git/go-git/v5/plumbing"            // https://pkg.go.dev/github.com/go-git/go-git/v5
	ggitobject "github.com/go-git/go-git/v5/plumbing/object"       // https://pkg.go.dev/github.com/go-git/go-git/v5
	ggithttp "github.com/go-git/go-git/v5/plumbing/transport/http" // https://pkg.go.dev/github.com/go-git/go-git/v5
	log "github.com/sirupsen/logrus"                               // https://pkg.go.dev/github.com/sirupsen/logrus

	errs "github.com/mooltiverse/nyx/modules/go/errors"
	gitent "github.com/mooltiverse/nyx/modules/go/nyx/entities/git"
)

/*
A local repository implementation that encapsulates the backing go-git (https://pkg.go.dev/github.com/go-git/go-git/v5) library.
*/
type goGitRepository struct {
	// The private instance of the underlying Git object.
	repository *ggit.Repository
}

/*
Builds the instance using the given backing object.
*/
func newGoGitRepository(repository *ggit.Repository) (goGitRepository, error) {
	if repository == nil {
		return goGitRepository{}, &errs.NilPointerError{Message: fmt.Sprintf("nil pointer '%s'", "repository")}
	}

	gitRepository := goGitRepository{}
	gitRepository.repository = repository
	return gitRepository, nil
}

/*
Returns a new basic authentication method object using the given user name and password.

Returns nil if both the given credentials are nil.

  - user the user name to use when credentials are required. It may be nil.
    When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
  - password the password to use when credentials are required. It may be nil.
    When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
*/
func getBasicAuth(user *string, password *string) *ggithttp.BasicAuth {
	if user == nil && password == nil {
		return nil
	} else if user != nil && password == nil {
		return &ggithttp.BasicAuth{Username: *user}
	} else if user == nil && password != nil {
		return &ggithttp.BasicAuth{Password: *password}
	} else {
		return &ggithttp.BasicAuth{Username: *user, Password: *password}
	}
}

/*
Returns a repository instance working in the given directory after cloning from the given URI.

Arguments are as follows:

- directory the directory where the repository has to be cloned. It is created if it doesn't exist.
- uri the URI of the remote repository to clone.

Errors can be:

- NilPointerError if any of the given objects is nil
- IllegalArgumentError if the given object is illegal for some reason, like referring to an illegal repository
- GitError in case the operation fails for some reason, including when authentication fails
*/
func clone(directory *string, uri *string) (goGitRepository, error) {
	return cloneWithCredentials(directory, uri, nil, nil)
}

/*
Returns a repository instance working in the given directory after cloning from the given URI.

Arguments are as follows:

  - directory the directory where the repository has to be cloned. It is created if it doesn't exist.
  - uri the URI of the remote repository to clone.
  - user the user name to use when credentials are required. If this and password are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
  - password the password to use when credentials are required. If this and user are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.

Errors can be:

- NilPointerError if any of the given objects is nil
- IllegalArgumentError if the given object is illegal for some reason, like referring to an illegal repository
- GitError in case the operation fails for some reason, including when authentication fails
*/
func cloneWithCredentials(directory *string, uri *string, user *string, password *string) (goGitRepository, error) {
	if directory == nil {
		return goGitRepository{}, &errs.NilPointerError{Message: "can't clone a repository instance with a null directory"}
	}
	if uri == nil {
		return goGitRepository{}, &errs.NilPointerError{Message: "can't clone a repository instance with a null URI"}
	}
	if "" == strings.TrimSpace(*directory) {
		return goGitRepository{}, &errs.IllegalArgumentError{Message: "can't create a repository instance with a blank directory"}
	}
	if "" == strings.TrimSpace(*uri) {
		return goGitRepository{}, &errs.IllegalArgumentError{Message: "can't create a repository instance with a blank URI"}
	}

	log.Debugf("cloning repository in directory '%s' from URI '%s'", *directory, *uri)

	repository, err := ggit.PlainClone(*directory, false, &ggit.CloneOptions{URL: *uri, Auth: getBasicAuth(user, password)})
	if err != nil {
		return goGitRepository{}, &errs.GitError{Message: fmt.Sprintf("unable to clone the '%s' repository into '%s'", *uri, *directory), Cause: err}
	}

	return newGoGitRepository(repository)
}

/*
Returns a repository instance working in the given directory.

Arguments are as follows:

- directory the directory where the repository is.

Errors can be:

- IllegalArgumentError if the given object is illegal for some reason, like referring to an illegal repository
- IOError in case of any I/O issue accessing the repository
*/
func open(directory string) (goGitRepository, error) {
	if "" == strings.TrimSpace(directory) {
		return goGitRepository{}, &errs.IllegalArgumentError{Message: "can't create a repository instance with a blank directory"}
	}
	repository, err := ggit.PlainOpen(directory)
	if err != nil {
		return goGitRepository{}, &errs.IllegalArgumentError{Message: fmt.Sprintf("unable to open Git repository in directory '%s'", directory), Cause: err}
	}
	return newGoGitRepository(repository)
}

/*
Resolves the commit with the given id using the repository object and returns it as a typed object.

This method is an utility wrapper around CommitObject which never returns
nil and throws GitError if the identifier cannot be resolved or any other error occurs.

Arguments are as follows:

- id the commit identifier to resolve. It must be a long or abbreviated SHA-1 but not nil.

Errors can be:

- GitError in case the given identifier cannot be resolved or any other issue is encountered
*/
func (r goGitRepository) parseCommit(id string) (ggitobject.Commit, error) {
	log.Tracef("parsing commit '%s'", id)
	commit, err := r.repository.CommitObject(ggitplumbing.NewHash(id))
	if err != nil {
		return ggitobject.Commit{}, &errs.GitError{Message: fmt.Sprintf("the '%s' commit identifier cannot be resolved as there is no such commit.", id), Cause: err}
	}
	return *commit, nil
}

/*
Resolves the object with the given id in the repository.

This method is an utility wrapper around ResolveRevision which never returns
nil and returns GitError if the identifier cannot be resolved or any other error occurs.

Arguments are as follows:

  - id the object identifier to resolve. It can't be nil. If it's a SHA-1 it can be long or abbreviated.
    For allowed values see ResolveRevision

Errors can be:

- GitError in case the given identifier cannot be resolved or any other issue is encountered
*/
func (r goGitRepository) resolve(id string) (ggitplumbing.Hash, error) {
	log.Tracef("resolving '%s'", id)

	rev, err := r.repository.ResolveRevision(ggitplumbing.Revision(id))
	if err != nil {
		return ggitplumbing.Hash{}, &errs.GitError{Message: fmt.Sprintf("the '%s' identifier cannot be resolved", id), Cause: err}
	}
	if rev == nil {
		if "HEAD" == id {
			log.Warnf("Repository identifier '%s' cannot be resolved. This means that the repository has just been initialized and has no commits yet or the repository is in a 'detached HEAD' state. See the documentation to fix this.", "HEAD")
		}
		return ggitplumbing.Hash{}, &errs.GitError{Message: fmt.Sprintf("Identifier '%s' cannot be resolved", id)}
	} else {
		return ggitplumbing.NewHash(rev.String()), nil
	}
}

/*
Arguments are as follows:

- paths the file patterns of the contents to add to stage. Cannot be nil or empty. The path "." represents
all files in the working area so with that you can add all locally changed files.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to add paths.
*/
func (r goGitRepository) Add(paths []string) error {
	log.Debugf("adding contents to repository staging area")
	if paths == nil || len(paths) == 0 {
		return &errs.GitError{Message: fmt.Sprintf("cannot stage a nil or empty set of paths")}
	}

	worktree, err := r.repository.Worktree()
	if err != nil {
		return &errs.GitError{Message: fmt.Sprintf("an error occurred when getting the current worktree for the repository"), Cause: err}
	}
	for _, path := range paths {
		err := worktree.AddWithOptions(&ggit.AddOptions{All: false, Path: "", Glob: path})
		if err != nil {
			return &errs.GitError{Message: fmt.Sprintf("an error occurred when trying to add paths to the staging area"), Cause: err}
		}
	}

	return nil
}

/*
Commits changes to the repository. Files to commit must be staged separately using Add.

- message the commit message. Cannot be nil.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to commit.
*/
func (r goGitRepository) CommitWithMessage(message *string) (gitent.Commit, error) {
	return r.CommitWithMessageAndIdentities(message, nil, nil)
}

/*
Commits changes to the repository. Files to commit must be staged separately using Add.

Arguments are as follows:

- message the commit message. Cannot be nil.
- author the object modelling the commit author informations. It may be nil, in which case the default
for the repository will be used
- committer the object modelling the committer informations. It may be nil, in which case the default
for the repository will be used

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to commit.
*/
func (r goGitRepository) CommitWithMessageAndIdentities(message *string, author *gitent.Identity, committer *gitent.Identity) (gitent.Commit, error) {
	log.Debugf("committing changes to repository")

	if message == nil {
		return gitent.Commit{}, &errs.GitError{Message: fmt.Sprintf("cannot commit with a nil message")}
	}

	worktree, err := r.repository.Worktree()
	if err != nil {
		return gitent.Commit{}, &errs.GitError{Message: fmt.Sprintf("an error occurred when getting the current worktree for the repository"), Cause: err}
	}
	var gAuthor *ggitobject.Signature = nil
	var gCommitter *ggitobject.Signature = nil
	if author != nil {
		gAuthor = &ggitobject.Signature{Name: author.Name, Email: author.Email}
	}
	if committer != nil {
		gCommitter = &ggitobject.Signature{Name: committer.Name, Email: committer.Email}
	}
	commitHash, err := worktree.Commit(*message, &ggit.CommitOptions{All: false, Author: gAuthor, Committer: gCommitter})
	if err != nil {
		return gitent.Commit{}, &errs.GitError{Message: fmt.Sprintf("an error occurred when trying to commit"), Cause: err}
	}
	commit, err := r.repository.CommitObject(commitHash)
	if err != nil {
		return gitent.Commit{}, &errs.GitError{Message: fmt.Sprintf("an error occurred when retrieving the commit that has been created"), Cause: err}
	}
	return CommitFrom(*commit, []gitent.Tag{}), nil
}

/*
Adds the given files to the staging area and commits changes to the repository. This method is a shorthand
for Add and CommitWithMessage.

Arguments are as follows:

  - paths the file patterns of the contents to add to stage. Cannot be nil or empty. The path "." represents
    all files in the working area so with that you can add all locally changed files.
  - message the commit message. Cannot be nil.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to commit.
*/
func (r goGitRepository) CommitPathsWithMessage(paths []string, message *string) (gitent.Commit, error) {
	return r.CommitPathsWithMessageAndIdentities(paths, message, nil, nil)
}

/*
Adds the given files to the staging area and commits changes to the repository. This method is a shorthand
for Add and CommitWithMessageAndIdentities.

Arguments are as follows:

  - paths the file patterns of the contents to add to stage. Cannot be nil or empty. The path "." represents
    all files in the working area so with that you can add all locally changed files.
  - message the commit message. Cannot be nil.
  - author the object modelling the commit author informations. It may be nil, in which case the default
    for the repository will be used
  - committer the object modelling the committer informations. It may be nil, in which case the default
    for the repository will be used

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to commit.
*/
func (r goGitRepository) CommitPathsWithMessageAndIdentities(paths []string, message *string, author *gitent.Identity, committer *gitent.Identity) (gitent.Commit, error) {
	err := r.Add(paths)
	if err != nil {
		return gitent.Commit{}, &errs.GitError{Message: fmt.Sprintf("an error occurred while staging contents to the repository"), Cause: err}
	}
	return r.CommitWithMessageAndIdentities(message, author, committer)
}

/*
Returns a set of abjects representing all the tags for the given commit.

Arguments are as follows:

- commit the SHA-1 identifier of the commit to get the tags for. It can be a full or abbreviated SHA-1.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository.
*/
func (r goGitRepository) GetCommitTags(commit string) ([]gitent.Tag, error) {
	log.Debugf("retrieving tags for commit '%s'", commit)
	var res []gitent.Tag
	tagsIterator, err := r.repository.Tags()
	if err != nil {
		return nil, &errs.GitError{Message: fmt.Sprintf("cannot list repository tags"), Cause: err}
	}
	if err := tagsIterator.ForEach(func(ref *ggitplumbing.Reference) error {
		// in order to check if the tag has this commit as target we first need to figure out if it's annotated or lightweight
		tagObject, err := r.repository.TagObject(ref.Hash())
		switch err {
		case nil:
			// it's an annotated tag
			if strings.HasPrefix(tagObject.Target.String(), commit) {
				res = append(res, TagFrom(r.repository, *ref))
			}
		case ggitplumbing.ErrObjectNotFound:
			// it's a lightweight tag
			if strings.HasPrefix(ref.Hash().String(), commit) {
				res = append(res, TagFrom(r.repository, *ref))
			}
		default:
			// Some other error occurred
			return &errs.GitError{Message: fmt.Sprintf("error while listing repository tags"), Cause: err}
		}
		return nil
	}); err != nil {
		return nil, &errs.GitError{Message: fmt.Sprintf("error while listing repository tags"), Cause: err}
	}
	return res, nil
}

/*
Returns the name of the current branch or a commit SHA-1 if the repository is in the detached head state.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or is in the 'detached HEAD' state.
*/
func (r goGitRepository) GetCurrentBranch() (string, error) {
	ref, err := r.repository.Head()
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("unable to resolve reference to HEAD"), Cause: err}
	}

	// also strip the leading "refs/heads/" from the reference name
	return strings.Replace(ref.Name().String(), "refs/heads/", "", 1), nil
}

/*
Returns the SHA-1 identifier of the last commit in the current branch.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or is in the 'detached HEAD' state.
*/
func (r goGitRepository) GetLatestCommit() (string, error) {
	ref, err := r.repository.Head()
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("unable to resolve reference to HEAD"), Cause: err}
	}
	commitSHA := ref.Hash().String()
	log.Debugf("repository latest commit in HEAD branch is '%s'", commitSHA)
	return commitSHA, nil
}

/*
Returns the SHA-1 identifier of the first commit in the repository (the only commit with no parents).

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or is in the 'detached HEAD' state.
*/
func (r goGitRepository) GetRootCommit() (string, error) {
	ref, err := r.repository.Head()
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("unable to resolve reference to HEAD"), Cause: err}
	}
	// the Log method doesn't let us follow the firt parent, so we need to go through all commits and stop at the end
	commit, err := r.parseCommit(ref.Hash().String())
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("an error occurred while walking the commit history at commit '%s'", ref.Hash().String()), Cause: err}
	}
	for len(commit.ParentHashes) > 0 {
		c, err := r.repository.CommitObject(commit.ParentHashes[0]) // always follow the first parent, ignore others, if any
		if err != nil {
			return "", &errs.GitError{Message: fmt.Sprintf("an error occurred while walking the commit history at commit '%s'", ref.Hash().String()), Cause: err}
		}
		commit = *c
	}
	commitSHA := commit.Hash.String()
	log.Debugf("repository latest commit in HEAD branch is '%s'", commitSHA)
	return commitSHA, nil
}

/*
Returns the names of configured remote repositories.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or is in the 'detached HEAD' state.
*/
func (r goGitRepository) GetRemoteNames() ([]string, error) {
	log.Debugf("retrieving repository remote names")
	remotes, err := r.repository.Remotes()
	if err != nil {
		return nil, &errs.GitError{Message: fmt.Sprintf("unable to get the repository remotes"), Cause: err}
	}
	remoteNames := make([]string, len(remotes))
	for i, rmt := range remotes {
		remoteNames[i] = rmt.Config().Name
	}

	log.Debugf("repository remote names are '%v'", remoteNames)
	return remoteNames, nil
}

/*
Returns true if the repository is clean, which is when no differences exist between the working tree, the index,
and the current HEAD.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or is in the 'detached HEAD' state.
*/
func (r goGitRepository) IsClean() (bool, error) {
	log.Debugf("checking repository clean status")
	wt, err := r.repository.Worktree()
	if err != nil {
		return false, &errs.GitError{Message: fmt.Sprintf("unable to get the repository worktree"), Cause: err}
	}
	status, err := wt.Status()
	if err != nil {
		return false, &errs.GitError{Message: fmt.Sprintf("unable to get the repository worktree status"), Cause: err}
	}
	return status.IsClean(), nil
}

/*
Pushes local changes in the current branch to the default remote origin.

# Returns the local name of the remotes that has been pushed

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) Push() (string, error) {
	return r.PushWithCredentials(nil, nil)
}

/*
Pushes local changes in the current branch to the default remote origin.

Returns the local name of the remotes that has been pushed.

Arguments are as follows:

  - user the user name to create when credentials are required. If this and password are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
  - password the password to create when credentials are required. If this and user are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) PushWithCredentials(user *string, password *string) (string, error) {
	s := DEFAULT_REMOTE_NAME
	return r.PushToRemoteWithCredentials(&s, user, password)
}

/*
Pushes local changes in the current branch to the given remote.

Returns the local name of the remotes that has been pushed.

Arguments are as follows:

- remote the name of the remote to push to. If nil or empty the default remote name (origin) is used.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) PushToRemote(remote *string) (string, error) {
	return r.PushToRemoteWithCredentials(remote, nil, nil)
}

/*
Pushes local changes in the current branch to the default remote origin.

Returns the local name of the remotes that has been pushed.

Arguments are as follows:

  - remote the name of the remote to push to. If nil or empty the default remote name (origin) is used.
  - user the user name to create when credentials are required. If this and password are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
  - password the password to create when credentials are required. If this and user are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) PushToRemoteWithCredentials(remote *string, user *string, password *string) (string, error) {
	remoteString := ""
	if remote != nil {
		remoteString = *remote
	}
	log.Debugf("pushing changes to remote repository '%s'", remoteString)

	// get the current branch name
	ref, err := r.repository.Head()
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("unable to resolve reference to HEAD"), Cause: err}
	}
	currentBranchRef := ref.Name()
	// the refspec is in the localBranch:remoteBranch form, and we assume they both have the same name here
	branchRefSpec := ggitconfig.RefSpec(currentBranchRef + ":" + currentBranchRef)
	tagsRefSpec := ggitconfig.RefSpec("refs/tags/*:refs/tags/*") // this is required to also push tags

	err = r.repository.Push(&ggit.PushOptions{RemoteName: remoteString, RefSpecs: []ggitconfig.RefSpec{branchRefSpec, tagsRefSpec}, Auth: getBasicAuth(user, password)})
	if err != nil {
		return "", &errs.GitError{Message: fmt.Sprintf("an error occurred when trying to push"), Cause: err}
	}
	return remoteString, nil
}

/*
Pushes local changes in the current branch to the given remotes.

Returns a collection with the local names of remotes that have been pushed.

Arguments are as follows:

- remotes the names of remotes to push to. If nil or empty the default remote name (origin) is used.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) PushToRemotes(remotes []string) ([]string, error) {
	return r.PushToRemotesWithCredentials(remotes, nil, nil)
}

/*
Pushes local changes in the current branch to the given remotes.

Returns a collection with the local names of remotes that have been pushed.

Arguments are as follows:

  - remotes remotes the names of remotes to push to. If nil or empty the default remote name (origin) is used.
  - user the user name to create when credentials are required. If this and password are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.
  - password the password to create when credentials are required. If this and user are both nil
    then no credentials is used. When using single token authentication (i.e. OAuth or Personal Access Tokens)
    this value may be the token or something other than a token, depending on the remote provider.

Errors can be:

- GitError in case some problem is encountered with the underlying Git repository, preventing to push.
*/
func (r goGitRepository) PushToRemotesWithCredentials(remotes []string, user *string, password *string) ([]string, error) {
	log.Debugf("pushing changes to '%d' remote repositories", len(remotes))
	var res []string
	for _, remote := range remotes {
		r, err := r.PushToRemoteWithCredentials(&remote, user, password)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

/*
Tags the latest commit in the current branch with a tag with the given name. The resulting tag is lightweight.

Returns the object modelling the new tag that was created. Never nil.

Arguments are as follows:

- name the name of the tag. Cannot be nil

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, preventing to tag
    (i.e. when the tag name is nil or there is already a tag with the given name in the repository).
*/
func (r goGitRepository) Tag(name *string) (gitent.Tag, error) {
	return r.TagWithMessage(name, nil)
}

/*
Tags the latest commit in the current branch with a tag with the given name and optional message.

Returns the object modelling the new tag that was created. Never nil.

Arguments are as follows:

  - name the name of the tag. Cannot be nil
  - message the optional tag message. If nil the new tag will be lightweight, otherwise it will be an
    annotated tag

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, preventing to tag
    (i.e. when the tag name is nil or there is already a tag with the given name in the repository).
*/
func (r goGitRepository) TagWithMessage(name *string, message *string) (gitent.Tag, error) {
	return r.TagWithMessageAndIdentity(name, message, nil)
}

/*
Tags the latest commit in the current branch with a tag with the given name and optional message using the optional
tagger identity.

Returns the object modelling the new tag that was created. Never nil.

Arguments are as follows:

  - name the name of the tag. Cannot be nil
  - message the optional tag message. If nil the new tag will be lightweight, otherwise it will be an
    annotated tag
  - tagger the optional identity of the tagger. If nil Git defaults are used. If message is nil this is ignored.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, preventing to tag
    (i.e. when the tag name is nil or there is already a tag with the given name in the repository).
*/
func (r goGitRepository) TagWithMessageAndIdentity(name *string, message *string, tagger *gitent.Identity) (gitent.Tag, error) {
	return r.TagCommitWithMessageAndIdentity(nil, name, message, tagger)
}

/*
Tags the object represented by the given SHA-1 with a tag with the given name and optional message using the optional
tagger identity.

Returns the object modelling the new tag that was created. Never nil.

Arguments are as follows:

  - target the SHA-1 identifier of the object to tag. If nil the latest commit in the current branch is tagged.
  - name the name of the tag. Cannot be nil
  - message the optional tag message. If nil the new tag will be lightweight, otherwise it will be an
    annotated tag
  - tagger the optional identity of the tagger. If nil Git defaults are used. If message is nil this is ignored.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, preventing to tag
    (i.e. when the tag name is nil or there is already a tag with the given name in the repository).
*/
func (r goGitRepository) TagCommitWithMessageAndIdentity(target *string, name *string, message *string, tagger *gitent.Identity) (gitent.Tag, error) {
	if name == nil {
		return gitent.Tag{}, &errs.GitError{Message: fmt.Sprintf("tag name cannot be nil")}
	}

	log.Debugf("tagging as '%s'", *name)
	var createTagOptions *ggit.CreateTagOptions = nil
	if message != nil {
		var gTagger *ggitobject.Signature = nil
		if tagger != nil {
			gTagger = &ggitobject.Signature{Name: tagger.Name, Email: tagger.Email}
		}
		// create an annotated tag, pass a CreateTagOptions
		// when the message is nil we create a lightweight tag so CreateTagOptions needs to be nil
		createTagOptions = &ggit.CreateTagOptions{Tagger: gTagger, Message: *message}
	}
	var targetHash ggitplumbing.Hash
	if target == nil {
		commitSHA, err := r.GetLatestCommit()
		if err != nil {
			return gitent.Tag{}, &errs.GitError{Message: fmt.Sprintf("unable to get the latest commit (HEAD)"), Cause: err}
		}
		targetHash = ggitplumbing.NewHash(commitSHA)
	} else {
		targetHash = ggitplumbing.NewHash(*target)
	}
	ref, err := r.repository.CreateTag(*name, targetHash, createTagOptions)

	if err != nil {
		return gitent.Tag{}, &errs.GitError{Message: fmt.Sprintf("unable to create Git tag"), Cause: err}
	}
	return TagFrom(r.repository, *ref), nil
}

/*
Browse the repository commit history using the given visitor to inspect each commit. Commits are
evaluated in Git's natural order, from the most recent to oldest.

Arguments are as follows:

  - start the optional SHA-1 id of the commit to start from. If nil the latest commit in the
    current branch (HEAD) is used. This can be a long or abbreviated SHA-1. If this commit cannot be
    resolved within the repository a GitError is thrown.
  - end the optional SHA-1 id of the commit to end with, included. If nil the repository root
    commit is used (until the given visitor returns false). If this commit is not reachable
    from the start it will be ignored. This can be a long or abbreviated SHA-1. If this commit cannot be resolved
    within the repository a GitError is thrown.
  - visit the visitor function that will receive commit data to evaluate. If nil this method takes no action.
    The function isits a single commit and receives all of the commit simplified fields. Returns true
    to keep browsing next commits or false to stop.

Errors can be:

  - GitError in case some problem is encountered with the underlying Git repository, including when
    the repository has no commits yet or a given commit identifier cannot be resolved.
*/
func (r goGitRepository) WalkHistory(start *string, end *string, visit func(commit gitent.Commit) bool) error {
	if visit == nil {
		return nil
	}
	startString := "not defined"
	if start != nil {
		startString = *start
	}
	endString := "not defined"
	if end != nil {
		endString = *end
	}
	log.Debugf("walking commit history. Start commit boundary is '%s'. End commit boundary is '%s'", startString, endString)
	log.Debugf("upon merge commits only the first parent is considered.")

	var commit *ggitobject.Commit
	if start == nil {
		startHash, err := r.GetLatestCommit()
		if err != nil {
			return err
		}
		c, err := r.parseCommit(startHash)
		commit = &c
		if err != nil {
			return err
		}
	} else {
		c, err := r.parseCommit(*start)
		commit = &c
		if err != nil {
			return err
		}
	}
	log.Tracef("start boundary resolved to commit '%s'", commit.Hash.String())

	if end != nil {
		// make sure it can be resolved
		c, err := r.parseCommit(*end)
		endCommit := &c
		if err != nil {
			return err
		}
		log.Tracef("end boundary resolved to commit '%s'", endCommit.Hash.String())
	}

	for commit != nil {
		log.Tracef("visiting commit '%s'", commit.Hash.String())

		tags, err := r.GetCommitTags(commit.Hash.String())
		if err != nil {
			return err
		}
		visitorContinues := visit(CommitFrom(*commit, tags))

		if !visitorContinues {
			log.Debugf("commit history walk interrupted by visitor")
			break
		} else if end != nil && strings.HasPrefix(commit.Hash.String(), *end) {
			log.Debugf("commit history walk reached the end boundary '%s'", *end)
			break
		} else if len(commit.ParentHashes) == 0 {
			commit = nil
			log.Debugf("commit history walk reached the end")
			break
		} else {
			commit, err = r.repository.CommitObject(commit.ParentHashes[0]) // follow the first parent upon merge commits
			if err != nil {
				return &errs.GitError{Message: fmt.Sprintf("an error occurred while walking through commits"), Cause: err}
			}
		}
	}
	return nil
}