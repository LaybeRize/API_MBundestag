# MBundestag (a political simulation)
<img style="float: left" width="159px" src="https://github.com/LaybeRize/API_MBundestag/blob/main/public/MBundestagLogo.png?raw=true" alt="MBundestag Logo">

[![Release](https://img.shields.io/badge/Version-0.10.0-blue)](https://img.shields.io)

The goal of this project is to create a webserver application, that 
provides all the necessary pages and functionalities to run a political
simulation.

With this goal in mind, the project tries to model its user input options and interactions, as well as the 
moderation options after real world political situations. The simulation also strifes to be able to model 
all and every type of government, so that any type of ideology and all types of resistance against the system
can be at least to a certain degree accurately modeled.

## Details on the project's goal

The most important accomplishment is the actual real world like 
documentation of politics in a controlled environment. There are 
three things that have to be possible. The hidden organisation 
of events, actions and discussions. A possibility to vote on matters and
a way to document changes and the history of a document itself, apart
from being able to publish news in a unified way and making contracts.

Another feature that I would like to see implemented is a way to communicate
in a twitter-like way but with more characters to give an alternative
to the streamlined political discussion described above.

## General Implementations needed

* Login-system
* User-self info and management (contracts, password changes, list of organisation memberships and titles owned)
* Admin-handling (contract verification, organisation management, user management, title management)
* Newspaper (creation, approval, viewing)
* Bills/Posts, Discussions and Votes (creation and management)

# Requirements

The simulation uses Postgresql to save all permanent and semi-permanent data. The project provides a finished 
dockerfile and docker-composer.yaml to build the project and deploy it.

## Postgresql

The data will be saved in a postgresql database, that is linked to the server by envoirnment variables. 
These are as follows:
````text
DB_NAME=yourDBName
DB_PASSWORD=yourDBPassword
DB_USER=yourDBUser
DOCKER=yourIPAdressForPostgres
ADRESS=yourIPAdressForAccessToThePage
INIT_NAME=TheRootAccountDisplayname
INIT_USERNAME=TheRootAccountUsername
INIT_PASSWORD=TheRootAccountPassword
````

DOCKER and ADRESS should be formatted like this, when used with the provided dockerfile and composer:


````text
DB_NAME=yourDBName
DB_PASSWORD=yourDBPassword
DB_USER=yourDBUser
DOCKER=db
ADRESS=0.0.0.0
INIT_NAME=TheRootAccountDisplayname
INIT_USERNAME=TheRootAccountUsername
INIT_PASSWORD=TheRootAccountPassword
````

all other parameter are free to be configured by the hoster. 
## Accounts

The user's account will be created by the admins. For the possibility 
to write newspapers from a different personality without logging out there
should be a system in place, to assign the user at least one press account.
Press accounts should be normal accounts in all but access.

Accounts have flairs, they are the conglomerate of all titles and membership tokens of the account.
For press accounts this can also mean a specific set affiliation, set by the admins.
Flairs are added to all posts of a user, but the flairs shown for a specific post do 
not change, when the user flair changes. That keeps the historical circumstances of the post itself correct.

## Organisations

Every Post, Discussion or Vote must be published under the name of an organisation.
organisations can be private, public or secret. Public organisation's posts, discussion and votes
are for all to see, at all times. Private organisation's posts are always public, their availability of
votes and discussion however can be limited to a specific group of people, but must not.
Secret organisation's votes and discussion can only be seen by member of the organisation.
They can not create posts. (In a way they can put out decrees with discussions 
that have no one that is allowed to comment)

## Newspaper

Newspaper should serve as the informational backbone of the simulation. 
Here are one can find the general ongoings in the world of the simulation.
These can inform events and decisions triggered by members of the simulation alike.
Everyone can publish an article, but the final authority of what makes it to the public's 
eyes are the admins. They have to manage the information flow through the newspaper.
Especially no-canon events should be caught and deleted by them.

## Contracts or Letters

Interally called Letters, because they serve multiple puropses. The system provides Users with a 
way to inform of specific terms or contract them to specific terms themselves. Letters are only visable 
to people added by the creator themselves. They can only be viewed by admins, when given the UUID of the letter.
Specifically for court cases where a contract was breached or an agreement, that contains illegal actions, 
the letter can be used by the prosecution, if one of the contracts is willed to give the UUID to the Admins to 
verify it contents. Letters are also a way for the moderation to reject articles while informing the user why and
what should be changed to make the article acceptable. Also Moderation managed events or contracts with moderation-simulated 
firms and organisations can be handled in this way. It can also be way to ensure the validity of an interview, by 
contracting the interviewed person to allow them to publish the interview in an article.

# Plan

All milestones and functionalities are here documented, for progress overview purpose.

## Config parameter

* [x] Can set name and password for standard head admin account
* [x] Can set port for application

## Accounts

* [x] Accounts implemented for database
* [x] Account can be created
* [x] Press Account can be linked to specific User-Account
* [x] Account can be suspended
* [x] When changing an account you specifiy if the account should be removed from all organisations it resides in and also if it should lose all titles 
* [x] Account flairs and access level can be changed
* [x] Accounts can be grouped by owner
* [x] Accounts are ranked into user, media admin, admin and head admin
* [x] media admins can approve and deny newspaper entries
* [x] admins can manage organisations and titles
* [x] head admins can create accounts and suspend user
* [x] there is a standard head admin account created
* [x] Accounts fully functional

## Newspaper

* [x] Articles and Publications implemented for database
* [x] Articles can be submitted
* [x] Articles can be rejected before publication (giving the user back a letter that can be extended with a reason for the rejection)
* [x] Articles are being grouped in publications
* [x] Breaking-News can be published separate from the normal newspaper
* [x] All news are pending for publication until ratified by the moderation 
* [x] Admins can publish a newspapers and breaking news
* [x] Publishing a normal newspaper automatically groups any new entrances into the next newspaper
* [x] All accounts can manage, edit, create and delete newspaper based on their access level

## Titles

* [x] Titles have a main- and a subgroup they belong to
* [x] They can be displayed in a hierarchy grouped and sorted by their main group, subgroup and their name, respectively.
* [x] There exists a view of the Organigramm of all titles.
* [x] Titles also list how they appear in flairs and who are currently holding them
* [x] A title can have more than one owner
* [x] The title flair is automatically added to the list of flairs of the user
* [x] Titles are fully functional editable by admins and visible to all user

### Note:

* [x] Flair value in database must be either unique or null
* [x] Flairs are not allowed to have commas

## Contracts

* [x] Contracts (in database called letters) implemented for database
* [x] Contracts consist of a body and a title. As well as a list of accounts signing.
* [x] User can send a contract to any amount of people involved. They are officially called letters.
* [x] User can specify that a contract only shows any signature when all user have signed. If any specified user refuses to sign, the contract will be rejected
* [x] User can sign or refuse a contract. If not all user must sign, the signature will be viewable by any person part of the contract immediately.
* [x] User can also specifiy, that the letter is view only
* [x] The letter function is also available to admins sending mod mails, imitating character and person, not yet present
* [x] Contracts can be viewed by admins given the UUID, the user can see on their contract.
* [x] Admin View for Modmessages

## Self-view

* [x] User can look at their own profile
* [x] They must be able to see, what press accounts they own
* [x] They must be able to change their accounts password
* [x] They must be able to view their contracts based on their accounts
* [x] They must be able to see what flairs they have and what titles they own as well as what organisations they are part of

## Zwitscher

* [x] Zwitscher implemented for database
* [x] For short, social-media-like communication there should be a platform
* [x] That platform should allow anyone to post with any of their accounts a zwitscher consisting of up to 500 characters.
* [x] To every zwitscher everyone can comment.
* [x] Comments are zwitscher itself.

## Organisations

* [x] Organisations have a main- and a subgroup they belong to
* [x] They can be either public, private, secret or hidden.
* [x] They can be displayed in a hierarchy grouped and sorted by their main group, subgroup and their name, respectively.
* [x] Secret organisations can not publish posts
* [x] Secret organisation's votes and discussion can not be set to public
* [x] Private organisation get to choose who can view discussions and votes
* [x] Public/Private Organisations have the public trait automatically selected for discussions and votes.
* [x] There exists a view of the Organigramm of all non-hidden organisations.
* [x] There exists a table view of all hidden organisations.
* [x] Members are either admins or normal users.
* [x] Admin-Members can modify the user-Members of an organisation
* [x] User get their assigned flair for their memberships
* [x] Flairs are only assignable to Organisations that are public or private
* [x] Private and Public organisations publish how is Adminstrator and User as part of the organisation and their assigned flair, if it exists
* [x] Accounts can manage Organisations based on their assigned role and position

### Note:

* [x] Flair value in database must be either unique or null
* [x] Flairs are not allowed to have commas

## Posts 

* [x] Posts implemented for database
* [x] Posts can be created by admins of organisations only.
* [x] Posts have a history that can be extended by organisation admins.
* [x] Posts have a title and a subtitle as well as a body.
* [x] Posts can be viewed by anyone

````html
<i class="bi bi-file-text"></i>
````

## Discussions

* [x] Discussion implemented for database
* [x] Discussions have title, subtitle and body.
* [x] Public discussion can specify who is allowed to comment.
* [x] Private discussion can specify who is allowed to see, but not comment.
* [x] Discussions for private and secrete organisations are automatically private.
* [x] Discussions are streamlined. A comment can only be made and is added at the end of the current discussion. No sub-comments are possible.
* [x] Discussion comments can be deleted by admins, displaying only to the user, that there once was a comment. The content will be held in the database for bookkeeping.

````html
<i class="bi bi-chat-right-text"></i>
````

## Votes

* [x] Votes implemented for database
* [x] Votes consist of a title, subtitle, body and the polls itself.
* [x] Every poll consists of the viewable data while the poll is going on and after.
* [x] Every private vote can specify who can view the vote, aside from the people that can participate.
* [x] Every vote must have a finish date. All polls close simultaneously after that date.

````html
<i class="bi bi-archive"></i>
````

## Discord API

* [ ] Connect with Discord via Token given in the envoirment variable
* [ ] Stream newspaper to dedicated channel
* [ ] Stream posts to channel
* [ ] make abo for notifications on votes and discussions

## Additional Features for the feature

* [ ] Implement a chat function to make interviews possible on the side with only your press account revealed (maybe over discord?)
* [ ] Logging all admin action
* [ ] make actions from the log revertable
* [ ] Excel like tables to log infos for admins

# Design

The design of the website should look serious. In a way an old, simple look (clear sharp edges, mono color) and a 
document like font would be preferable.

## Designing with the provided templating engine

* [x] implemented a markdown to html converter + styling options
* [x] top down templating
* [x] easy way to prevent html code duplication with ways to customize the duplicate
* [x] possibility to leave out attributes which will be filled with an empty space in the template

Assuming we need a constant backdrop for the content like this.

```html
<!-- backdrop.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <title></title>
</head>
<body>
<!-- Content goes here -->
</body>
```

We can modify the file like this with the typical go template syntax.

```html
<!-- backdrop.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <title></title>
</head>
<body>
{{block "content" .}}
{{end}}
</body>
```

These kind of files must be stored in the layouts folder in templates.
For easy usage in a page just write the content like you would normally define it 
with go templates like this and add a comment with the layout name.


```html
<!-- page.html -->
{{define "content"}}
<p>Test paragraph Element</p>
{{end}}

{{/* use backdrop */}}
```

These files must be stored under pages folder in the templates.
Any extra includes are then defined in html files in the includes folder and can be 
used like this:
```html
<!-- page.html -->
{{define "content"}}
<p>Test paragraph Element</p>
{{template "includes/test"}}
<!-- the file is just called test.html -->
{{end}}

{{/* use backdrop */}}
<!-- makes the template aware of what backdrop the content should be -->
<!-- inserted into. In this case the backdrop.html -->
```

If you want to use the new templating replace engine written you have to create a html in the elements-folder.
Something like this.

```html
<!-- elements.html -->
The file first starts parsing after a html comment is made like this
<!-- test: -->
<p>This will be inserted instead of the tag</p>
```

Or if you want to use changeable parameter like this

```html
<!-- elements.html -->
The file first starts parsing after a html comment is made like this
<!-- test-param: param1, param2 -->
<p>This will be inserted instead #param1# of the tag and also add #param2#</p>
<p>#content#</p>
```
The examples would be used just like this in the other html files:

```html
<!-- random html code -->
<div>
    <test/>
    <test-param param1="this goes here" param2="this goes there">
        Here goes what will be copied to #content# 
    </test-param>
</div>
<!-- random html code -->
```

The above example will be saved in cache like this:

```html
<!-- random html code -->
<div>
    <p>This will be inserted instead of the tag</p>
    <p>This will be inserted instead this goes here of the tag and also add this goes there</p>
    <p>Here goes what will be copied to #content#</p>
</div>
<!-- random html code -->
```

The order of the variables is not relevant. The template will now compile if one or more parameter are omitted. 
If the parameter is omitted, the variable will be replaced with an empty string.

Additional Information: If you want to use the quotes in the attribute you have to do it like this

```html
<!-- random html code -->
<div>
    <test/>
    <test-param param1="this goes here &quot;Quotes go hard&quot;" param2="this goes there">
        Here goes what will be copied to #content#
    </test-param>
</div>
<!-- random html code -->
```

And the result should then be

```html
<!-- random html code -->
<div>
    <p>This will be inserted instead of the tag</p>
    <p>This will be inserted instead this goes here &quot;Quotes go hard&quot; of the tag and also add this goes there</p>
    <p>Here goes what will be copied to #content#</p>
</div>
<!-- random html code -->
```
### Markdown Converter

the ``markdown.html`` in the includes folder is a styling guide for what the html view of the 
user's markdown will look like. The given attributes are extracted from the ``markdown.html``
and then applied, when the ``helper.CreateHTML`` is called on the markdown-string.
## Webinterfaces

The following pages are needed.

* [x] Login or Homepage
* [x] Account management page
* [x] Organisation management page
* [x] Title management page
* [x] Press management page
* [x] Article view page
* [x] Publication view page
* [x] Organisation organigramm page
* [x] title organigramm page
* [x] admin contract view page
* [x] user self-view/password edit page
* [x] contract list view page
* [x] zwitscher start page
* [x] zwitscher post view page
* [x] post view page
* [x] discussion view page
* [ ] vote view page
* [x] list views for post discussion and vote (implementation details pending)
* [ ] Impressum my beloved


# Specific implementation details
Any details not specified in the checklist are subject to change and discussion 
of the development and design team. They will be listed here, if a satisfactory solution 
is found.

# Mentions

css and fonts from https://github.com/twbs/icons copied, to keep it local.
package htmlWrapper is an extended version of https://github.com/Xeoncross/got.
Every image used in the project (aside from the logo which is the offical german eagle) 
is from https://unsplash.com/.