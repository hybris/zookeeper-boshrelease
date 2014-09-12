# BOSH Release for Zookeeper

## Table of Contents
* [Usage](#usage)
* [Releases](#releases)
* [Backup](#backup)
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
  <tr>
    <td>2014-06-24</td>
    <td>zookeeper-hybris-2.yml</td>
    <td>2</td>
    <td>-</td>
    <td>Adding offset value for index to support cluster in multizones</td>
  </tr>
  <tr>
    <td>2014-06-24</td>
    <td>zookeeper-hybris-3.yml</td>
    <td>3</td>
    <td>-</td>
    <td>small fix for cluster support</td>
  </tr>
</table>


## <a name="backup"></a>Zookeeper-Backup

To have the Zookeeper backed up regularly we have created a Go program that does backups regularly.
You will find the source code at `src/zk-backyp`.

### How does it work?
Zookeeper is an in-memory, distributed, non-sharded, consistent, fault-tolerant database. That means, on every node you have the same information. When one node fails its workload will be transferred to another node. You will get a higher latency for a moment but no query should get lost. All data is held in memory and is written to a transaction log on disk. From time to time a complete snapshot of the in-memory data is written to a snapshot file. The next transaction log starts from this snapshot.
We use this characteristic to implement a backup mechanism:
1. Shut down one Zookeeper node
2. tar czf snapshots and their incrementing transaction logs
3. Restart the Zookeeper node
4. Upload the backup .tar.gz to S3
5. Delete the backup .tar.gz from local disk

####The Backup process runs only on ONE Zookeeper node!

To deploy the ZK-Backyp programm on an instance, add the zk-backyp template to the deployment manifest:

    - name: kafka_zookeeper_prod_us-east-1c
      instances: 1
      templates:
        - name: zookeeper
          release: zookeeper-hybris
        - name: zk-backyp
          release: zookeeper-hybris
      ...

Further you have to set some properties to configure the backup name (built from prefix + datestring + .tar.gz), backup interval in seconds (default 86400 seconds = 1 day) as well as the S3 target.

    properties:
      zk_backyp:
        prefix: zookeeper-backup-eu-west1-staged
        interval: 86400
        s3_bucket: zookeeper-backup.cf2.hybris.com
        s3_endpoint: https://s3-eu-west-1.amazonaws.com
        aws_access_key: AJHKJUAFH6FLFJKH5
        aws_secret_key: dhihfAFDSdh4j3FSFDAFh35j4lFAD/jklhl43fe