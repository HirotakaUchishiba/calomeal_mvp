# terraform/modules/vpc/main.tf

# 1. VPC本体を作成
resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr
  tags = {
    Name = "${var.project_name}-vpc"
  }
}

# 2. インターネットへの出口となるゲートウェイを作成
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "${var.project_name}-igw"
  }
}

# 3. パブリックサブネット（ALBなどを配置）を作成
resource "aws_subnet" "public" {
  count             = length(var.public_subnets_cidr)
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.public_subnets_cidr[count.index]
  availability_zone = "${var.region}${element(["a", "c"], count.index)}"
  map_public_ip_on_launch = true # パブリックIPを自動で割り当てる
  tags = {
    Name = "${var.project_name}-public-subnet-${count.index + 1}"
  }
}

# 4. プライベートサブネット（DBやアプリを配置）を作成
resource "aws_subnet" "private" {
  count             = length(var.private_subnets_cidr)
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.private_subnets_cidr[count.index]
  availability_zone = "${var.region}${element(["a", "c"], count.index)}"
  tags = {
    Name = "${var.project_name}-private-subnet-${count.index + 1}"
  }
}

# 5. パブリックサブネット用のルートテーブル（交通整理ルール）を作成
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id
  route {
    cidr_block = "0.0.0.0/0" # インターネット向けの通信は
    gateway_id = aws_internet_gateway.main.id # インターネットゲートウェイに向ける
  }
  tags = {
    Name = "${var.project_name}-public-rt"
  }
}

# 6. パブリックサブネットとルートテーブルを紐付け
resource "aws_route_table_association" "public" {
  count          = length(aws_subnet.public)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}