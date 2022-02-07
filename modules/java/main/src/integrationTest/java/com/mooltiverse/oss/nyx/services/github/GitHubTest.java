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
package com.mooltiverse.oss.nyx.services.github;

import static org.junit.jupiter.api.Assertions.*;

import java.util.Map;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;

import com.mooltiverse.oss.nyx.git.Scenario;
import com.mooltiverse.oss.nyx.git.Script;
import com.mooltiverse.oss.nyx.git.util.RandomUtil;
import com.mooltiverse.oss.nyx.services.Release;
import com.mooltiverse.oss.nyx.services.Service;

@DisplayName("GitHub")
public class GitHubTest {
    @Nested
    @DisplayName("Instance")
    class InstanceTest {
        @Test
        public void exceptionWithNullOptions()
            throws Exception {
            assertThrows(NullPointerException.class, () -> GitHub.instance(null));
        }

        @Test
        public void instanceWithEmptyOptions()
            throws Exception {
            GitHub service = GitHub.instance(Map.<String,String>of());
            assertNotNull(service);
        }
    }

    @Nested
    @DisplayName("Supports")
    class SupportsTest {
        @ParameterizedTest(name = "GitHub.instance().supports(''{0}'') == true")
        @EnumSource(Service.Feature.class)
        public void supportAnyFeature(Service.Feature feature)
            throws Exception {
            if (Service.Feature.GIT_HOSTING.equals(feature) || Service.Feature.RELEASES.equals(feature) || Service.Feature.USERS.equals(feature))
                assertTrue(GitHub.instance(Map.<String,String>of()).supports(feature));
            else assertFalse(GitHub.instance(Map.<String,String>of()).supports(feature));
        }
    }

    @Nested
    @DisplayName("User Service")
    class UserServiceTest {
        @Test
        @DisplayName("GitHub.instance().getAuthenticatedUser() throws SecurityException with no authentication token")
        public void exceptionUsingGetAuthenticatedUserWithNoAuthenticationToken()
            throws Exception {
            assertThrows(com.mooltiverse.oss.nyx.services.SecurityException.class, () -> GitHub.instance(Map.<String,String>of()).getAuthenticatedUser());
        }

        @Test
        @DisplayName("GitHub.instance().getAuthenticatedUser()")
        public void getAuthenticatedUser()
            throws Exception {
            // the 'gitHubTestUserToken' system property is set by the build script, which in turn reads it from an environment variable
            GitHubUser user = GitHub.instance(Map.<String,String>of(GitHub.AUTHENTICATION_TOKEN_OPTION_NAME, System.getProperty("gitHubTestUserToken"))).getAuthenticatedUser();
            assertNotNull(user.getID());
            assertNotNull(user.getUserName());
            assertNotNull(user.geFullName());
            assertNotNull(user.getAttributes());
            assertFalse(user.getAttributes().isEmpty());
        }
    }

    @Nested
    @DisplayName("Hosting Service")
    class HostingServiceTest {
        @Test
        @DisplayName("GitHub.instance().createGitRepository() and GitHub.instance().deleteGitRepository()")
        public void createGitRepository()
            throws Exception {
            String randomID = RandomUtil.randomAlphabeticString(5);
            // the 'gitHubTestUserToken' system property is set by the build script, which in turn reads it from an environment variable
            GitHub gitHub = GitHub.instance(Map.<String,String>of(GitHub.AUTHENTICATION_TOKEN_OPTION_NAME, System.getProperty("gitHubTestUserToken")));
            GitHubUser user = gitHub.getAuthenticatedUser();
            GitHubRepository gitHubRepository = gitHub.createGitRepository(randomID, "Test repository "+randomID, true, false);
            
            assertEquals("main", gitHubRepository.getDefaultBranch());
            assertEquals("Test repository "+randomID, gitHubRepository.getDescription());
            assertEquals(randomID, gitHubRepository.getName());
            assertEquals(user.getUserName()+"/"+randomID, gitHubRepository.getFullName());
            assertEquals("https://github.com/"+user.getUserName()+"/"+randomID+".git", gitHubRepository.getHTTPURL());

            // if we delete too quickly we often get a 404 from the server so let's wait a short while
            Thread.sleep(2000);

            // now delete it
            gitHub.deleteGitRepository(randomID);
        }
    }

    @Nested
    @DisplayName("Release Service")
    class ReleaseServiceTest {
        @Test
        @DisplayName("GitHub.instance().createRelease() and GitHub.instance().getRelease()")
        public void createRelease()
            throws Exception {
            String randomID = RandomUtil.randomAlphabeticString(5);
            // the 'gitHubTestUserToken' system property is set by the build script, which in turn reads it from an environment variable
            GitHub gitHub = GitHub.instance(Map.<String,String>of(GitHub.AUTHENTICATION_TOKEN_OPTION_NAME, System.getProperty("gitHubTestUserToken")));
            GitHubUser user = gitHub.getAuthenticatedUser();
            GitHubRepository gitHubRepository = gitHub.createGitRepository(randomID, "Test repository "+randomID, false, true);

            assertNull(gitHub.getReleaseByTag(user.getUserName(), gitHubRepository.getName(), "1.0.0-alpha.1"));

            // if we clone too quickly next calls may fail
            Thread.sleep(2000);

            // when a token for user and password authentication for plain Git operations against a GitHub repository,
            // the user is the token and the password is the empty string
            Script script = Scenario.FIVE_BRANCH_UNMERGED_BUMPING_COLLAPSED.applyOnClone(gitHubRepository.getHTTPURL(), System.getProperty("gitHubTestUserToken"), "");
            script.push(System.getProperty("gitHubTestUserToken"), "");

            // publish the release
            Release release = gitHub.publishRelease(user.getUserName(), gitHubRepository.getName(), "Release 1.0.0-alpha.1", "1.0.0-alpha.1", "A test description for the release\non multiple lines\nlike these");
            assertEquals("Release 1.0.0-alpha.1", release.getTitle());
            assertEquals("1.0.0-alpha.1", release.getTag());

            // now delete it
            gitHub.deleteGitRepository(randomID);
        }
    }
}