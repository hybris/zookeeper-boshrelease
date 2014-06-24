# BOSH Release for Zookeeper

## Table of Contents
* [Usage](#usage)
* [Releases](#releases)
<br />
<br />


### <a name="usage"></a>Usage

To use this bosh release, first upload it to your bosh:

```
bosh target BOSH_HOST
bosh login
git clone ssh://git@stash.hybris.com:7999/idefix/bosh-release-zookeeper.git
cd bosh-release-zookeeper
bosh upload release releases/zookeeper-hybris-XXXX.yml
```

### <a name="releases"></a>Releases
<table>
  <tr>
    <th>Date</th>
    <th>Name</th>
    <th>Version</th>
    <th>Zookeeper</th>
    <th>Comment</th>
  </tr>
  <tr>
    <td>2014-05-08</td>
    <td>zoo.hybris-1.yml</td>
    <td>1</td>
    <td>-</td>
    <td>-</td>
  </tr>
  <tr>
    <td>2014-06-04</td>
    <td>hybris-zookeeper-3.yml</td>
    <td>3</td>
    <td>-</td>
    <td>-</td>
  </tr>
  <tr>
    <td>2014-06-05</td>
    <td>hybris-zookeeper-4.yml</td>
    <td>3</td>
    <td>-</td>
    <td>-</td>
  </tr>
  <tr>
    <td>2014-06-24</td>
    <td>zookeeper-hybris-1.yml</td>
    <td>1</td>
    <td>-</td>
    <td>NEW RELEASE NAME/VERSION - COMPLETE REDESIGN!</td>
  </tr>
</table>
