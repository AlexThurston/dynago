package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"io/ioutil"
	"math/rand"
	"strconv"
)

var config = aws.NewConfig().WithEndpoint("http://localhost:16666").WithRegion("share-local")
var db = dynamodb.New(config)
var images = make(map[string][]string)

func PrintTableNames() {
	listTableParams := &dynamodb.ListTablesInput{}
	resp, err := db.ListTables(listTableParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

func ScanImageTable() {
	params := &dynamodb.ScanInput{
		TableName:       aws.String("Image"),
		AttributesToGet: []*string{},
	}

	resp, err := db.Scan(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(resp)
	fmt.Println(*resp.Count)

}

func ScanImageTagTable() {
	params := &dynamodb.ScanInput{
		TableName: aws.String("ImageTag"),
		AttributesToGet: []*string{
			aws.String("Tag"),
			aws.String("ImageId"),
			aws.String("VoteCount"),
		},
	}

	resp, err := db.Scan(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(*resp.Count)
}

func AddImages() {
	for image, _ := range images {
		addItemParams := &dynamodb.PutItemInput{
			TableName: aws.String("Image"),
			Item: map[string]*dynamodb.AttributeValue{
				"Id": &dynamodb.AttributeValue{
					S: aws.String(image),
				},
				"DateAdded": &dynamodb.AttributeValue{
					S: aws.String("today"),
				},
				"VoteCount": &dynamodb.AttributeValue{
					N: aws.String(strconv.Itoa(rand.Intn(100))),
				},
			},
		}
		_, err := db.PutItem(addItemParams)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func AddImageTags() {
	for image, tags := range images {
		for _, tag := range tags {
			addItemParams := &dynamodb.PutItemInput{
				TableName: aws.String("ImageTag"),
				Item: map[string]*dynamodb.AttributeValue{
					"Tag": &dynamodb.AttributeValue{
						S: aws.String(tag),
					},
					"ImageId": &dynamodb.AttributeValue{
						S: aws.String(image),
					},
					"VoteCount": &dynamodb.AttributeValue{
						N: aws.String(strconv.Itoa(rand.Intn(100))),
					},
				},
			}
			_, err := db.PutItem(addItemParams)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func QueryImageTagTable() {
	params := &dynamodb.QueryInput{
		TableName: aws.String("ImageTag"),
		KeyConditions: map[string]*dynamodb.Condition{
			"Tag": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("Application Services"),
					},
				},
			},
		},
	}
	resp, err := db.Query(params)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}

func main() {
	images["android.png"] = []string{"SDKs & Tools", "Android"}
	images["appstream.png"] = []string{"Application Services", "Amazon AppStream"}
	images["cli.png"] = []string{"SDKs & Tools", "AWS CLI"}
	images["cloudformation.png"] = []string{"Deployment & Management", "AWS CloudFormation"}
	images["cloudfront.png"] = []string{"Storage & CDN", "Amazon CloudFront"}
	images["cloudsearch.png"] = []string{"Application Services", "Amazon CloudSearch"}
	images["cloudtrail.png"] = []string{"Deployment & Management", "AWS CloudTrail"}
	images["cloudwatch.png"] = []string{"Deployment & Management", "Amazon CloudWatch"}
	images["data-pipeline.png"] = []string{"Analytics", "AWS Data Pipeline"}
	images["direct-connect.png"] = []string{"Compute & Networking", "AWS Direct Connect"}
	images["dotnet.png"] = []string{"SDKs & Tools", ".NET"}
	images["dynamodb.png"] = []string{"Database", "Amazon DynamoDB"}
	images["ec2.png"] = []string{"Compute & Networking", "Amazon EC2"}
	images["eclipse.png"] = []string{"SDKs & Tools", "Eclipse"}
	images["elasticache.png"] = []string{"Database", "Amazon ElastiCache"}
	images["elastic-beanstalk.png"] = []string{"Deployment & Management", "AWS Elastic Beanstalk"}
	images["elb.png"] = []string{"Compute & Networking", "Elastic Load Balancing"}
	images["emr.png"] = []string{"Analytics", "Amazon EMR"}
	images["glacier.png"] = []string{"Storage & CDN", "Amazon Glacier"}
	images["iam.png"] = []string{"Deployment & Management", "AWS IAM"}
	images["ios.png"] = []string{"SDKs & Tools", "iOS"}
	images["java.png"] = []string{"SDKs & Tools", "Java"}
	images["kinesis.png"] = []string{"Analytics", "Amazon Kinesis"}
	images["nodejs.png"] = []string{"SDKs & Tools", "Node.js"}
	images["opsworks.png"] = []string{"Deployment & Management", "AWS OpsWorks"}
	images["php.png"] = []string{"SDKs & Tools", "PHP"}
	images["powershell.png"] = []string{"SDKs & Tools", "PowerShell"}
	images["python.png"] = []string{"SDKs & Tools", "Python"}
	images["rds.png"] = []string{"Database", "Amazon RDS"}
	images["redshift.png"] = []string{"Database", "Amazon Redshift"}
	images["route53.png"] = []string{"Compute & Networking", "Amazon Route 53"}
	images["ruby.png"] = []string{"SDKs & Tools", "Ruby"}
	images["s3.png"] = []string{"Storage & CDN", "Amazon S3"}
	images["ses.png"] = []string{"Application Services", "Amazon SES"}
	images["sns.png"] = []string{"Application Services", "Amazon SNS"}
	images["sqs.png"] = []string{"Application Services", "Amazon SQS"}
	images["storage-gateway.png"] = []string{"Storage & CDN", "Amazon Storage Gateway"}
	images["swf.png"] = []string{"Application Services", "Amazon SWF"}
	images["transcoding.png"] = []string{"Application Services", "Amazon Elastic Transcoder"}
	images["visual-studio.png"] = []string{"SDKs & Tools", "Visual Studio"}
	images["vpc.png"] = []string{"Compute & Networking", "Amazon VPC"}

	// Print all the table names
	PrintTableNames()

	// Create the image table
	imageJson, err := ioutil.ReadFile("./Image.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var imageTableCreateInput dynamodb.CreateTableInput
	err = json.Unmarshal(imageJson, &imageTableCreateInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.CreateTable(&imageTableCreateInput)
	if err != nil {
		fmt.Println(err)
	}

	// Create the image tag table
	imageTagJson, err := ioutil.ReadFile("./ImageTag.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	var imageTagTableCreateInput dynamodb.CreateTableInput
	err = json.Unmarshal(imageTagJson, &imageTagTableCreateInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.CreateTable(&imageTagTableCreateInput)
	if err != nil {
		fmt.Println(err)
	}

	AddImages()
	AddImageTags()
	ScanImageTable()
	ScanImageTagTable()
	QueryImageTagTable()

	// Add a row to the table
	// addItemParams := &dynamodb.PutItemInput{
	// 	TableName: aws.String("Users"),
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"Id": &dynamodb.AttributeValue{
	// 			N: aws.String("1"),
	// 		},
	// 		"FirstName": &dynamodb.AttributeValue{
	// 			S: aws.String("Alex"),
	// 		},
	// 		"LastName": &dynamodb.AttributeValue{
	// 			S: aws.String("Thurston"),
	// 		},
	// 		"Age": &dynamodb.AttributeValue{
	// 			S: aws.String("33"),
	// 		},
	// 	},
	// }

	// _, err = db.PutItem(addItemParams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Add another row to the table
	// addItemParams = &dynamodb.PutItemInput{
	// 	TableName: aws.String("Users"),
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"Id": &dynamodb.AttributeValue{
	// 			N: aws.String("2"),
	// 		},
	// 		"FirstName": &dynamodb.AttributeValue{
	// 			S: aws.String("Lauren"),
	// 		},
	// 		"LastName": &dynamodb.AttributeValue{
	// 			S: aws.String("Jolly"),
	// 		},
	// 		"Age": &dynamodb.AttributeValue{
	// 			S: aws.String("33"),
	// 		},
	// 	},
	// }

	// _, err = db.PutItem(addItemParams)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Read a row from the table
	// getItemParams := &dynamodb.GetItemInput{
	// 	TableName: aws.String("Users"),
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"Id": &dynamodb.AttributeValue{
	// 			N: aws.String("2"),
	// 		},
	// 	},
	// 	AttributesToGet: []*string{
	// 		aws.String("FirstName"),
	// 		aws.String("LastName"),
	// 		aws.String("Age"),
	// 	},
	// }
	// resp, getErr := db.GetItem(getItemParams)
	// if getErr != nil {
	// 	fmt.Println(getErr)
	// 	return
	// }
	// fmt.Println(resp)

	deleteImageTableParams := &dynamodb.DeleteTableInput{
		TableName: aws.String("Image"),
	}

	_, err = db.DeleteTable(deleteImageTableParams)
	if err != nil {
		fmt.Println(err)
	}

	deleteImageTagTableParams := &dynamodb.DeleteTableInput{
		TableName: aws.String("ImageTag"),
	}

	_, err = db.DeleteTable(deleteImageTagTableParams)
	if err != nil {
		fmt.Println(err)
	}

	PrintTableNames()

	fmt.Println("Done")
}
